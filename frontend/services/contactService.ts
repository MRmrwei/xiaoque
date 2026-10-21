import type { ContactInput } from '~/domain/types'

/**
 * Form submission boundary.
 *
 * Replace the body with `POST /api/v1/forms/contact/submissions` when the
 * backend is ready. Validation and UI state belong to useContactForm().
 */
export async function submitContact(_input: ContactInput): Promise<void> {
  await new Promise((resolve) => setTimeout(resolve, 450))
}
