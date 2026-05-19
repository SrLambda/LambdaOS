import pytest


@pytest.fixture(autouse=True)
def setup_textual():
    """Fixture que configura el entorno de Textual para testing.

    Se ejecuta automáticamente antes de cada test para asegurar que
    el entorno de Textual está correctamente inicializado y que no
    hay estado residual de tests anteriores.
    """
