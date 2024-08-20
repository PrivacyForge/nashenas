import nacl from 'tweetnacl'

function hexToUint8Array(hexString: string) {
  if (hexString.length % 2 !== 0) {
    throw new Error('Invalid hex string length')
  }

  const uint8Array = new Uint8Array(hexString.length / 2)

  for (let i = 0; i < hexString.length; i += 2) {
    uint8Array[i / 2] = parseInt(hexString.substr(i, 2), 16)
  }

  return uint8Array
}

function toHex(buffer: Uint8Array) {
  return Array.prototype.map
    .call(buffer, (x) => ('00' + x.toString(16)).slice(-2))
    .join('')
}

async function generateKeyPair() {
  const { publicKey, secretKey } = nacl.box.keyPair()
  return {
    publicKey: toHex(publicKey),
    privateKey: toHex(secretKey),
  }
}

async function createE2EPacket(
  recipientPublicKey: string,
  ownPrivateKey: string,
  message: string,
) {
  const decodedPublicKey = hexToUint8Array(recipientPublicKey)
  const decodedPrivateKey = hexToUint8Array(ownPrivateKey)
  const messageUint8Array = new TextEncoder().encode(message)
  const nonce = nacl.randomBytes(nacl.box.nonceLength)

  const encryptedMessage = nacl.box(
    messageUint8Array,
    nonce,
    new Uint8Array(decodedPublicKey),
    new Uint8Array(decodedPrivateKey),
  )
  const messagePackage = new Uint8Array(nonce.length + encryptedMessage.length)
  messagePackage.set(nonce)
  messagePackage.set(encryptedMessage, nonce.length)

  return toHex(messagePackage)
}

async function decryptE2EPacket(
  ownPrivateKey: string,
  recipientPublicKey: string,
  encryptedMessage: string,
) {
  const decodedPublicKey = hexToUint8Array(recipientPublicKey)
  const decodedPrivateKey = hexToUint8Array(ownPrivateKey)
  const decodedEncryptedMessage = hexToUint8Array(encryptedMessage)

  const receivedNonce = decodedEncryptedMessage.slice(0, nacl.box.nonceLength)
  const receivedEncryptedMessage = decodedEncryptedMessage.slice(
    nacl.box.nonceLength,
  )
  const decryptedMessage = nacl.box.open(
    receivedEncryptedMessage,
    receivedNonce,
    decodedPublicKey,
    decodedPrivateKey,
  )
  if (!decryptedMessage) return null
  return new TextDecoder().decode(decryptedMessage)
}

export { generateKeyPair, createE2EPacket, decryptE2EPacket }
