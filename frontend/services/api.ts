import axios from "axios"

const API_URL = "http://localhost:9090" // API Gateway URL

const api = axios.create({
  baseURL: API_URL,
  headers: {
    "Content-Type": "application/json",
  },
})

// Add a request interceptor to include the auth token in requests
api.interceptors.request.use(
  (config) => {
    console.log("API Request:", {
      method: config.method,
      url: config.url,
      data: config.data,
      headers: config.headers,
      baseURL: config.baseURL,
    });

    const token = localStorage.getItem("token")
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }
    return config
  },
  (error) => {
    console.error("API Request Error:", error);
    return Promise.reject(error)
  },
)

// Add a response interceptor to handle common errors
api.interceptors.response.use(
  (response) => {
    console.log("API Response:", {
      status: response.status,
      statusText: response.statusText,
      data: response.data,
      headers: response.headers,
      config: {
        method: response.config.method,
        url: response.config.url,
      },
    });
    return response
  },
  (error) => {
    console.error("API Response Error:", {
      message: error.message,
      code: error.code,
      response: error.response ? {
        status: error.response.status,
        statusText: error.response.statusText,
        data: error.response.data,
        headers: error.response.headers,
      } : 'No response',
      request: error.request ? 'Request exists' : 'No request',
      config: error.config ? {
        method: error.config.method,
        url: error.config.url,
        data: error.config.data,
      } : 'No config',
    });

    // Handle token expiration
    if (error.response && error.response.status === 401) {
      localStorage.removeItem("token")
      window.location.href = "/login"
    }

    return Promise.reject(error)
  },
)

export default api
