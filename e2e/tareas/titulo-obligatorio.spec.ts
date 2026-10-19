import { test, expect } from '../fixtures';

test.describe('Tareas', () => {
  test('el título es obligatorio', async ({ page }) => {
    // 1. Pulsa «Añadir» sin escribir nada
    await page.getByRole('button', { name: 'Añadir' }).click();

    // 2. Se muestra el error y no se crea ninguna tarea
    await expect(page.getByRole('alert')).toHaveText('El título es obligatorio.');
    await expect(page.getByRole('textbox', { name: 'Nueva tarea' })).toHaveAttribute('aria-invalid', 'true');
    await expect(page.locator('.task-title')).toHaveCount(0);
  });
});
