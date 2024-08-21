import * as AES from './AES'
import * as RSA from './RSA'

async function encrypt(
  message: string,
  sessionKey: string,
  publicKey: string,
) {
  const encryptedKey = await RSA.encrypt(sessionKey, publicKey)

  const encryptedMsg = await AES.encrypt(message, sessionKey)

  return `${encryptedKey}-${encryptedMsg}`
}

async function decrypt(message: string, privateKey: string, reverse: boolean) {
  let encryptedKey: string, encryptedMsg: string
  if (reverse) {
    ;[encryptedKey, encryptedMsg] = message.split('-').reverse()
  } else {
    ;[encryptedKey, encryptedMsg] = message.split('-')
  }

  const key = await RSA.decrypt(encryptedKey, privateKey)

  const msg = await AES.decrypt(encryptedMsg, key)

  return msg
}

export { encrypt, decrypt }
