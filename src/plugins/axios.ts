import Axios from 'axios'

alert(import.meta.env.VITE_BASE_URL)

const axios = Axios.create({
  baseURL: import.meta.env.VITE_BASE_URL,
})

axios.interceptors.request.use(async (config) => {
  config.headers['Authorization'] = window.Telegram.WebApp.initData
  return config
})

export default axios
