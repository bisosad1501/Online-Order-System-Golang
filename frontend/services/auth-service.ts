import api from "./api"

export const authService = {
  // Register a new user
  async register(userData: {
    email: string
    first_name: string
    last_name: string
    phone: string
    address: string
  }) {
    const response = await api.post("/users", userData)

    // Store the user ID in localStorage for convenience
    if (response.data && response.data.id) {
      localStorage.setItem("registeredUserId", response.data.id)
    }

    return response.data
  },

  // Login user with email only
  async login(credentials: {
    email: string
  }) {
    try {
      // Use the new login endpoint
      const response = await api.post("/auth/login", credentials);

      if (response.data.token && response.data.user) {
        // Store token and user ID in localStorage
        localStorage.setItem("token", response.data.token);
        localStorage.setItem("userId", response.data.user.id);

        return {
          token: response.data.token,
          user: response.data.user
        };
      } else {
        throw new Error("Login failed: Invalid response from server");
      }
    } catch (error) {
      console.error("Login error:", error);
      throw error;
    }
  },

  // Legacy login method using verify endpoint
  async verifyUser(credentials: {
    id?: string
    email: string
    address?: string
  }) {
    try {
      const response = await api.post("/users/verify", credentials);

      if (response.data.verified) {
        // Store token and user ID in localStorage
        if (response.data.token) {
          localStorage.setItem("token", response.data.token);
        }

        if (response.data.user && response.data.user.id) {
          localStorage.setItem("userId", response.data.user.id);
        }

        return response.data;
      } else {
        throw new Error(response.data.message || "Verification failed");
      }
    } catch (error) {
      console.error("Verification error:", error);
      throw error;
    }
  },

  // Get current user data
  async getCurrentUser() {
    // Get user ID from localStorage or JWT decode
    const token = localStorage.getItem("token");
    if (!token) {
      throw new Error("No token found");
    }

    // In a real app, we would decode the JWT to get the user ID
    // For now, we'll use a placeholder and assume the ID is stored in localStorage
    const userId = localStorage.getItem("userId");
    if (!userId) {
      throw new Error("No user ID found");
    }

    const response = await api.get(`/users/${userId}`);
    return response.data;
  },

  // Update user profile
  async updateProfile(
    userId: string,
    userData: {
      first_name?: string
      last_name?: string
      phone?: string
      address?: string
    },
  ) {
    const response = await api.put(`/users/${userId}`, userData)
    return response.data
  },
}
