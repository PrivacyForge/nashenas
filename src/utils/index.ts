interface IKeys {
  receivePublicKey: string
  receivePrivateKey: string
  sendPublicKey: string
  sendPrivateKey: string
}

function generateKeysTemplate(receiveKeys: string[], sendKeys: string[]) {
  return `${receiveKeys[0]}divide\n${receiveKeys[1]}separate\n${sendKeys[0]}divide\n${sendKeys[1]}`
}

function extractKeys(content: string) {
  const [receiveKeys, sendKeys] = content.split('separate')

  const receivePrivateKey = receiveKeys.split('divide')[0]
  const receivePublicKey = receiveKeys.split('divide')[1].slice(1)

  const sendPrivateKey = sendKeys.slice(1).split('divide')[0]
  const sendPublicKey = sendKeys.split('divide')[1].slice(1)

  return {
    receivePrivateKey,
    receivePublicKey,
    sendPrivateKey,
    sendPublicKey,
  }
}

function exportKeys() {
  const link = document.createElement('a')

  const receiveKeys = []
  receiveKeys[0] = localStorage.getItem('receive_private_key')!
  receiveKeys[1] = localStorage.getItem('receive_public_key')!

  const sendKeys = []
  sendKeys[0] = localStorage.getItem('send_private_key')!
  sendKeys[1] = localStorage.getItem('send_public_key')!

  const content = generateKeysTemplate(receiveKeys, sendKeys)

  const file = new Blob([content], { type: 'text/plain' })

  link.href = URL.createObjectURL(file)

  link.download = 'keys.txt'

  link.click()
  URL.revokeObjectURL(link.href)
}

function importKeysFromFile(file: File): Promise<IKeys> {
  return new Promise((resolve) => {
    const reader = new FileReader()
    reader.onload = (e) => {
      const rawData = e.target!.result as string

      const keys = extractKeys(rawData)

      resolve(keys)
    }
    reader.readAsText(file)
  })
}

function bufferToHex(buffer: ArrayBuffer) {
  return Array.from(new Uint8Array(buffer))
    .map((b) => b.toString(16).padStart(2, '0'))
    .join('')
}

export {
  generateKeysTemplate,
  extractKeys,
  exportKeys,
  importKeysFromFile,
  bufferToHex,
}
