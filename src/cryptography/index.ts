import * as AES from './AES'
import * as RSA from './RSA'

async function encrypt(
  message: string,
  destPubKey: string,
  srcPubKey?: string
) {
  const key = await AES.generateRandomKey()

  const destEncryptedKey = await RSA.encrypt(key, destPubKey)

  const encryptedMsg = await AES.encrypt(message, key)

  if (srcPubKey) {
    const srcEncryptedKey = await RSA.encrypt(key, srcPubKey)
    return `${destEncryptedKey}-${encryptedMsg}-${srcEncryptedKey}`
  }

  return `${destEncryptedKey}-${encryptedMsg}`
}

async function decrypt(message: string, privateKey: string) {
  const [encryptedKey, encryptedMsg] = message.split('-')

  const key = await RSA.decrypt(encryptedKey, privateKey)

  const msg = await AES.decrypt(encryptedMsg, key)

  return msg
}

export { encrypt, decrypt }
