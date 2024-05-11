import forge from 'node-forge'

async function generateKeyPair() {
  const keyPair = await forge.pki.rsa.generateKeyPair({ bits: 512, workers: 2 })
  const publicKey = forge.pki.publicKeyToPem(keyPair.publicKey)
  const privateKey = forge.pki.privateKeyToPem(keyPair.privateKey)

  return { privateKey, publicKey }
}

async function encrypt(message: string, publicKey: string) {
  const encryptedMsg = forge.pki.publicKeyFromPem(publicKey).encrypt(message)

  return forge.util.bytesToHex(encryptedMsg)
}

async function decrypt(message: string, privateKey: string) {
  return forge.pki
    .privateKeyFromPem(privateKey)
    .decrypt(forge.util.hexToBytes(message))
}

export { encrypt, decrypt, generateKeyPair }
