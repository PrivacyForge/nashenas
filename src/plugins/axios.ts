import Axios from 'axios'

const axios = Axios.create({
  baseURL: import.meta.env.VITE_BASE_URL,
})

axios.interceptors.request.use(async (config) => {
  config.headers['Authorization'] = `${localStorage.getItem('token')}`

  return config
})

export default axios
