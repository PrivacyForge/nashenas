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

export { generateKeysTemplate, extractKeys }
