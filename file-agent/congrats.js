/**
 * congrats.js
 * A script that decodes a ROT13-encoded message and prints it
 */

// The ROT13-encoded message
const encodedMessage = 'Pbatenghyngvbaf ba ohvyqvat n pbqr-rqvgvat ntrag!';

/**
 * Function to decode a ROT13-encoded string
 * @param {string} encodedStr - The ROT13-encoded string
 * @returns {string} The decoded string
 */
function rot13Decode(encodedStr) {
  return encodedStr.replace(/[a-zA-Z]/g, function(char) {
    // Get the character code
    const charCode = char.charCodeAt(0);
    
    // Handle uppercase letters (A-Z: 65-90)
    if (charCode >= 65 && charCode <= 90) {
      return String.fromCharCode(((charCode - 65 + 13) % 26) + 65);
    }
    
    // Handle lowercase letters (a-z: 97-122)
    if (charCode >= 97 && charCode <= 122) {
      return String.fromCharCode(((charCode - 97 + 13) % 26) + 97);
    }
    
    // Return non-alphabetic characters as is
    return char;
  });
}

// Decode and print the message
const decodedMessage = rot13Decode(encodedMessage);
console.log(decodedMessage);