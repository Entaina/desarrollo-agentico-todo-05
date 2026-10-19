import { test } from './fixtures';

// Punto de partida de los demás tests: la lista vacía y la página abierta. Sirve para
// explorar la aplicación con `playwright-cli attach` al escribir un test nuevo.
test('seed', async ({ page }) => {
  void page;
});
