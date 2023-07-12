// Needs to be adjusted for each structure of api error responses
export const getApiErrorMessage = (error?: { error?: string }) => {
  if (error?.error) {
    switch (error?.error) {
      case 'some_error_key':
        return 'Some error key value'
      default:
        error?.error
    }
  }
  return 'Internal server error has occurred. Please, contact system administrator.'
}

export const getHeaders = (token?: string): {} =>
  token
    ? {
        Authorization: `Bearer ${token}}`,
      }
    : {}

export const postHeaders = (token?: string): {} =>
  token
    ? {
        // 'Content-Type': 'application/json',
        Authorization: `Bearer ${token}}`,
      }
    : {
        // 'Content-Type': 'application/json',
      }
