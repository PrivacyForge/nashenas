import forge from 'node-forge'

const encoder = new TextEncoder()
const decoder = new TextDecoder()

function generateRandomKey(): string {
  const secretKey = forge.random.getBytesSync(16)
  const hexSecret = forge.util.bytesToHex(secretKey)

  return hexSecret
}

async function encrypt(message: string, key: string): Promise<string> {
  const data = encoder.encode(message)

  const cryptoKey = await crypto.subtle.importKey(
    'raw',
    encoder.encode(key),
    { name: 'AES-CBC' },
    false,
    ['encrypt']
  )

  const iv = crypto.getRandomValues(new Uint8Array(16))

  const encryptedData = await crypto.subtle.encrypt(
    { name: 'AES-CBC', iv },
    cryptoKey,
    data
  )

  const encryptedMessage = new Uint8Array(iv.length + encryptedData.byteLength)

  encryptedMessage.set(iv)
  encryptedMessage.set(new Uint8Array(encryptedData), iv.length)

  return Array.from(encryptedMessage)
    .map((byte) => ('0' + (byte & 0xff).toString(16)).slice(-2))
    .join('')
}

async function decrypt(
  encryptedMessage: string,
  secretKey: string
): Promise<string> {
  const encryptedData = new Uint8Array(
    encryptedMessage.match(/[\da-f]{2}/gi)!.map((hex) => parseInt(hex, 16))
  )

  const iv = encryptedData.slice(0, 16)

  const cryptoKey = await crypto.subtle.importKey(
    'raw',
    encoder.encode(secretKey),
    { name: 'AES-CBC' },
    false,
    ['decrypt']
  )

  const decryptedData = await crypto.subtle.decrypt(
    { name: 'AES-CBC', iv: iv },
    cryptoKey,
    encryptedData.slice(16)
  )

  return decoder.decode(decryptedData)
}

export { encrypt, decrypt, generateRandomKey }
