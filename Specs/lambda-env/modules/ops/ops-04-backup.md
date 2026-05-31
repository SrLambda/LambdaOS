# lambda-env: Backup Manager (ops-04)

## Intent

Módulo TUI para gestionar backups y snapshots: snapper, rsync, configuración de backup automático.

## Scope

- Snapper: listar snapshots, crear snapshot, rollback, comparar snapshots
- Rsync: configurar backup de directorios a destino local o remoto
- Configurar backup automático (cron/systemd timer)
- Ver espacio usado por snapshots

## Requirements

1. Snapper: listar snapshots con fecha, tipo (pre, post, manual), descripción
2. Crear snapshot manual con descripción
3. Rollback: seleccionar snapshot → restaurar
4. Comparar: ver diff entre dos snapshots
5. Rsync: configurar origen, destino, exclusions, schedule
6. Backend: `snapper`, `rsync`, systemd timers

## Scenarios

### Escenario 1: Crear snapshot antes de actualización
- Usuario abre Backup → Snapper → "Create snapshot"
- Descripción: "Before system update"
- Snapshot creado

### Escenario 2: Configurar backup de home a disco externo
- Abre Backup → Rsync → "New backup"
- Origen: `/home/lambda/`
- Destino: `/media/backup/lambda-home/`
- Excludes: `.cache/`, `node_modules/`, `.venv/`
- Schedule: daily at 2:00 AM
- Aplica → crea systemd timer

## Technical Notes

- Snapper requiere BTRFS (ya es el filesystem de LambdaOS)
- `snapper list` para listar snapshots
- `snapper create -d "description"` para crear
- `snapper rollback <number>` para rollback
- `rsync -avz --exclude=... src/ dest/` para backup
- Systemd timer para backups automáticos
- Agregar snapper a packages.x86_64 (mencionado en README pero no en packages)

## Dependencies

- `core/01-hub-plugin-system`
- `core/02-settings-schema`
- snapper (agregar a packages.x86_64)
- rsync (ya en packages.x86_64)
- Requiere root para snapper
