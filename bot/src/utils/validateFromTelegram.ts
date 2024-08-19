import { createHmac } from "node:crypto";

export const validateFromTelegram = (initData: string) => {
  function HMAC_SHA256(key: string | Buffer, secret: string) {
    return createHmac("sha256", key).update(secret);
  }
  function getCheckString(data: URLSearchParams) {
    const items: [k: string, v: string][] = [];
    // remove hash
    for (const [k, v] of data.entries()) if (k !== "hash") items.push([k, v]);

    return items
      .sort(([a], [b]) => a.localeCompare(b)) // sort keys
      .map(([k, v]) => `${k}=${v}`) // combine key-value pairs
      .join("\n");
  }

  const data = new URLSearchParams(initData);

  const data_check_string = getCheckString(data);
  const secret_key = HMAC_SHA256("WebAppData", process.env.BOT_TOKEN!).digest();
  const hash = HMAC_SHA256(secret_key, data_check_string).digest("hex");

  if (hash === data.get("hash")) return JSON.parse(data.get("user") || "{}");
  return false;
};
