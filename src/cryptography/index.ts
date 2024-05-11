import * as AES from './AES'
import * as RSA from './RSA'

async function encrypt(message: string, publicKey: string) {
  const key = await AES.generateRandomKey()

  const encryptedKey = await RSA.encrypt(key, publicKey)
  const encryptedMsg = await AES.encrypt(message, key)
  return `${encryptedKey}-${encryptedMsg}`
}

async function decrypt(message: string, privateKey: string) {
  const [encryptedKey, encryptedMsg] = message.split('-')

  const key = await RSA.decrypt(encryptedKey, privateKey)

  const msg = await AES.decrypt(encryptedMsg, key)

  return msg
}

export { encrypt, decrypt }
