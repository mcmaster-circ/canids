// Needs to be adjusted for each structure of api error responses
export const getApiErrorMessage = (error?: { message?: string }) => {
  if (error?.message) {
    return error?.message
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
