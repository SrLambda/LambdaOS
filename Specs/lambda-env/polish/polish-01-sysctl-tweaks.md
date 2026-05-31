# lambda-env: Sysctl Tweaks (polish-01)

## Intent

Aplicar optimizaciones de sysctl para mejorar rendimiento de red, memoria y filesystem en LambdaOS.

## Scope

- Optimizaciones de red: TCP buffers, BBR congestion control
- Optimizaciones de memoria: vm.swappiness, vm.vfs_cache_pressure
- Optimizaciones de filesystem: BTRFS tweaks
- Aplicar via sysctl.d

## Requirements

1. Archivo `/etc/sysctl.d/99-lambdaos.conf` con tweaks
2. Aplicado automáticamente al boot
3. No rompe funcionalidad existente
4. Documentado qué hace cada tweak

## Technical Notes

- Archivo: `airootfs/etc/sysctl.d/99-lambdaos.conf`
- Tweaks sugeridos:
  ```
  # Network
  net.core.default_qdisc = fq
  net.ipv4.tcp_congestion_control = bbr
  net.ipv4.tcp_fastopen = 3
  net.core.netdev_max_backlog = 65536

  # Memory
  vm.swappiness = 10
  vm.vfs_cache_pressure = 50
  vm.min_free_kbytes = 65536

  # BTRFS
  vm.dirty_ratio = 10
  vm.dirty_background_ratio = 5
  ```
- Verificar compatibilidad con kernel actual
- No aplicar tweaks que requieran módulos no cargados

## Dependencies

- Ninguno

## Verification

- `sysctl -p /etc/sysctl.d/99-lambdaos.conf` → aplica sin errores
- `sysctl net.ipv4.tcp_congestion_control` → muestra "bbr"
- `sysctl vm.swappiness` → muestra "10"
