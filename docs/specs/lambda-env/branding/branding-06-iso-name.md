# lambda-env: ISO Name & Publisher (branding-06)

## Intent

Personalizar el nombre de la ISO, publisher, y aplicación en `profiledef.sh` para que refleje la identidad de LambdaOS.

## Scope

- Cambiar `iso_name` de `lambda-os` a `LambdaOS` (o formato consistente)
- Cambiar `iso_publisher` a identidad oficial
- Cambiar `iso_application` a descripción correcta
- Configurar `iso_label` correctamente

## Requirements

1. ISO generada tiene nombre identificable
2. Publisher muestra identidad de LambdaOS
3. Application describe el propósito

## Technical Notes

- Archivo: `profiledef.sh`
- Valores actuales:
  - `iso_name="lambda-os"`
  - `iso_publisher="SrLambda <https://github.com/SrLambda>"`
  - `iso_application="Lambda OS Live/Rescue DVD"`
- Propuesta:
  - `iso_name="LambdaOS"`
  - `iso_publisher="LambdaOS Project <https://lambdaos.dev>"` (o URL que corresponda)
  - `iso_application="LambdaOS — The TUI-First Linux Distribution"`
- `iso_label` debe ser máximo 32 caracteres para compatibilidad ISO9660

## Dependencies

- Ninguno

## Verification

- ISO generada: `LambdaOS-<version>-x86_64.iso`
- `isoinfo -d -i LambdaOS-*.iso` → muestra publisher y application correctos
