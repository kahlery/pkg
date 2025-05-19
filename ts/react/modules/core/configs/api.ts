import axios from "axios"

import camelcaseKeys from "camelcase-keys"
import snakecaseKeys from "snakecase-keys"

export const Api = axios.create({
    baseURL: import.meta.env.VITE_API_URL,
})

Api.interceptors.request.use((config) => {
    if (config.data) {
        config.data = snakecaseKeys(config.data, { deep: true })
    }
})
