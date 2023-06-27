// Formatting messages recieved from server for display (using buefy snackbar)
export function format(data) {
  var jsonData = JSON.parse(data[1])
  var output = ""
  if (data[0] === false) {
    output = "Error: " + jsonData.message
  } else {
    output = jsonData.message
  }
  return output
}
