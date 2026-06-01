# Spec: dotfiles-module

## Intent

TUI module for GNU Stow operations: list stowed/unstowed modules, stow/unstow individual modules, detect file conflicts via checksums, and backup current configs to the dotfiles repo.

## Requirements

### Requirement 1: Module State Listing

The system SHALL list all dotfile modules in `~/dotfiles/` with their stowed/unstowed state. A module is considered stowed if its symlinks exist in the target directories.

#### Scenario: List modules with mixed state

- GIVEN `~/dotfiles/` contains directories: `nvim`, `qtile`, `kitty`
- AND `nvim` and `qtile` are stowed, `kitty` is not
- WHEN the module action is `list`
- THEN the output contains: `nvim: stowed`, `qtile: stowed`, `kitty: unstowed`

#### Scenario: Empty dotfiles directory

- GIVEN `~/dotfiles/` does not exist or is empty
- WHEN the module action is `list`
- THEN the module returns `{"status":"warning","code":"NO_MODULES","message":"No dotfiles modules found"}`

#### Scenario: Module directory without expected files

- GIVEN `~/dotfiles/kitty/` exists but contains no `.config/kitty/` subdirectory
- WHEN the module checks stow state
- THEN the module is listed as `unstowed`

### Requirement 2: Stow Operation

The system SHALL execute `stow <module>` from the `~/dotfiles/` directory to create symlinks. Before stowing, the system SHALL check for conflicts with existing files.

#### Scenario: Stow unstowed module

- GIVEN `kitty` module is unstowed
- WHEN the user selects "Stow" for `kitty`
- THEN the module executes `stow kitty` in `~/dotfiles/`
- AND symlinks are created in `~/.config/kitty/`
- AND the module returns `{"status":"ok","action":"stow","data":{"module":"kitty"}}`

#### Scenario: Stow already stowed module is idempotent

- GIVEN `nvim` module is already stowed
- WHEN the user selects "Stow" for `nvim`
- THEN the module executes `stow nvim` (no-op for existing symlinks)
- AND returns `{"status":"ok","action":"stow","data":{"module":"nvim","already_stowed":true}}`

#### Scenario: Stow with conflicts blocks operation

- GIVEN `~/.config/kitty/kitty.conf` exists and differs from `~/dotfiles/kitty/.config/kitty/kitty.conf`
- WHEN the user selects "Stow" for `kitty`
- THEN the module detects the conflict
- AND returns `{"status":"warning","code":"CONFLICT_DETECTED","message":"kitty.conf already exists and differs"}`
- AND stow is NOT executed

### Requirement 3: Unstow Operation

The system SHALL execute `stow -D <module>` to remove symlinks created by a previous stow operation.

#### Scenario: Unstow stowed module

- GIVEN `qtile` module is stowed
- WHEN the user selects "Unstow" for `qtile`
- THEN the module executes `stow -D qtile` in `~/dotfiles/`
- AND symlinks in `~/.config/qtile/` are removed
- AND the module returns `{"status":"ok","action":"unstow","data":{"module":"qtile"}}`

#### Scenario: Unstow already unstowed module

- GIVEN `kitty` module is not stowed
- WHEN the user selects "Unstow" for `kitty`
- THEN the module returns `{"status":"ok","action":"unstow","data":{"module":"kitty","already_unstowed":true}}`
- AND no stow command is executed

### Requirement 4: Conflict Detection via Checksums

The system SHALL detect conflicts by comparing file checksums between the dotfiles repo and the home directory. A conflict exists when a file exists in both locations with different content.

#### Scenario: No conflicts detected

- GIVEN all files in `~/dotfiles/nvim/` match their counterparts in `~/.config/nvim/`
- WHEN the module action is `check_conflicts`
- THEN the module returns `{"status":"ok","action":"check_conflicts","data":{"conflicts":[]}}`

#### Scenario: Conflict detected via checksum mismatch

- GIVEN `~/dotfiles/kitty/.config/kitty/kitty.conf` has checksum `abc123`
- AND `~/.config/kitty/kitty.conf` has checksum `def456`
- WHEN the module action is `check_conflicts`
- THEN the module returns `{"status":"ok","action":"check_conflicts","data":{"conflicts":[{"file":"kitty.conf","repo_checksum":"abc123","home_checksum":"def456"}]}}`

#### Scenario: File exists only in home (not a conflict)

- GIVEN `~/.config/nvim/init.lua` exists but `~/dotfiles/nvim/` has no `init.lua`
- WHEN the module action is `check_conflicts`
- THEN this file is NOT listed as a conflict (only differing files are conflicts)

### Requirement 5: Config Backup

The system SHALL backup current config files from the home directory to the dotfiles repo, preserving the directory structure. Before overwriting, the system SHALL check if the target file differs.

#### Scenario: Backup single module configs

- GIVEN `~/.config/kitty/kitty.conf` exists and differs from `~/dotfiles/kitty/.config/kitty/kitty.conf`
- WHEN the user selects "Backup" for `kitty`
- THEN the module copies `~/.config/kitty/` contents to `~/dotfiles/kitty/.config/kitty/`
- AND returns `{"status":"ok","action":"backup","data":{"module":"kitty","files_backed_up":1}}`

#### Scenario: Backup with no changes

- GIVEN all files in `~/.config/nvim/` match `~/dotfiles/nvim/`
- WHEN the user selects "Backup" for `nvim`
- THEN the module returns `{"status":"ok","action":"backup","data":{"module":"nvim","files_backed_up":0,"message":"No changes to backup"}}`

#### Scenario: Backup creates directory structure

- GIVEN `~/dotfiles/new-module/` exists but has no `.config/` subdirectory
- AND `~/.config/new-module/settings.yaml` exists
- WHEN the user selects "Backup" for `new-module`
- THEN the module creates `~/dotfiles/new-module/.config/new-module/`
- AND copies `settings.yaml` into it

## Technical Details

- Go package: `src/lambda-env/internal/modules/dotfiles/`
- Dotfiles root: `~/dotfiles/`
- Stow command: `stow <module>` / `stow -D <module>` (executed from dotfiles root)
- Checksum algorithm: SHA-256 via `crypto/sha256`
- Module manifest category: `ops`
- Module manifest dependencies: `["stow"]`

## Dependencies

- `core/01-hub-plugin-system` — module discovery and execution
- `core/02-settings-schema` — settings.json read/write
- GNU Stow installed on target system

## Verification Steps

```bash
# 1. Module compiles
cd src/lambda-env && go build ./internal/modules/dotfiles/...

# 2. Unit tests pass
cd src/lambda-env && go test ./internal/modules/dotfiles/... -v -cover

# 3. Stow/unstow works with test fixtures
mkdir -p /tmp/test-dotfiles/kitty/.config/kitty
echo "terminal = kitty" > /tmp/test-dotfiles/kitty/.config/kitty/kitty.conf
HOME=/tmp/test-home DOTFILES=/tmp/test-dotfiles stow -t /tmp/test-home -d /tmp/test-dotfiles kitty
# Verify: /tmp/test-home/.config/kitty/kitty.conf is a symlink

# 4. Conflict detection via checksums
# Create differing files in home vs dotfiles
# Run module check_conflicts action
# Expect: conflict listed with different checksums

# 5. Backup copies files correctly
# Modify file in home directory
# Run module backup action
# Expect: dotfiles repo file matches home directory
```
