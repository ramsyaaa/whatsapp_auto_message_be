// Dashboard JavaScript Logic

// Global variables
let token;
let messageChart, broadcastChart, hourlyChart;
let currentDays = 7;

// Initialize the dashboard when the DOM is loaded
document.addEventListener("DOMContentLoaded", function () {
  initializeDashboard();
});

// Main initialization function
function initializeDashboard() {
  // Check if user is logged in
  token = localStorage.getItem("admin_token");
  if (!token) {
    window.location.href = "/admin";
    return;
  }

  // Set user name
  const user = JSON.parse(localStorage.getItem("admin_user") || "{}");
  document.getElementById("user-name").textContent = user.name || "Admin User";

  // Initialize mobile sidebar
  initializeMobileSidebar();

  // Initialize logout functionality
  initializeLogout();

  // Initialize WhatsApp connection status
  initializeWhatsAppStatus();

  // Initialize period buttons
  initializePeriodButtons();

  // Initialize date picker
  initializeDatePicker();

  // Load all dashboard data
  loadAllDashboardData();

  // Set up auto-refresh
  setupAutoRefresh();
}

// Mobile sidebar functionality
function initializeMobileSidebar() {
  const mobileSidebar = document.getElementById("mobile-sidebar");
  const openSidebarBtn = document.getElementById("open-sidebar-btn");
  const closeSidebarBtn = document.getElementById("close-sidebar-btn");
  const sidebarBackdrop = document.getElementById("sidebar-backdrop");

  openSidebarBtn.addEventListener("click", () => {
    mobileSidebar.style.display = "flex";
  });

  closeSidebarBtn.addEventListener("click", () => {
    mobileSidebar.style.display = "none";
  });

  sidebarBackdrop.addEventListener("click", () => {
    mobileSidebar.style.display = "none";
  });
}

// Logout functionality
function initializeLogout() {
  const logoutBtn = document.getElementById("logout-btn");
  const mobileLogoutBtn = document.getElementById("mobile-logout-btn");

  function handleLogout() {
    localStorage.removeItem("admin_token");
    localStorage.removeItem("admin_user");
    window.location.href = "/admin";
  }

  logoutBtn.addEventListener("click", handleLogout);
  mobileLogoutBtn.addEventListener("click", handleLogout);
}

// WhatsApp connection status
function initializeWhatsAppStatus() {
  const statusIndicator = document.getElementById("status-indicator");
  const statusText = document.getElementById("status-text");
  const scanQrBtn = document.getElementById("scan-qr-btn");
  const logoutWhatsappBtn = document.getElementById("logout-whatsapp-btn");
  const qrContainer = document.getElementById("qr-container");
  const qrCode = document.getElementById("qr-code");

  // Check WhatsApp connection status
  checkWhatsAppStatus();

  // Scan QR code
  scanQrBtn.addEventListener("click", async () => {
    try {
      qrContainer.classList.remove("hidden");
      qrCode.src = "/api/v1/authenticate";
    } catch (error) {
      console.error("Error getting QR code:", error);
    }
  });

  // Logout WhatsApp
  logoutWhatsappBtn.addEventListener("click", async () => {
    try {
      const response = await fetch("/api/v1/log-out", {
        method: "POST",
        headers: {
          Authorization: `Bearer ${token}`,
        },
      });

      if (response.ok) {
        checkWhatsAppStatus();
      }
    } catch (error) {
      console.error("Error logging out from WhatsApp:", error);
    }
  });

  // Check WhatsApp status function
  async function checkWhatsAppStatus() {
    try {
      const response = await fetch("/api/v1/authenticate", {
        headers: {
          Authorization: `Bearer ${token}`,
        },
      });

      const data = await response.json();

      if (data.meta.status === "OK") {
        // Connected
        statusIndicator.classList.remove("bg-gray-300", "bg-red-500");
        statusIndicator.classList.add("bg-green-500");
        statusText.textContent = "Connected to WhatsApp";
        scanQrBtn.classList.add("hidden");
        logoutWhatsappBtn.classList.remove("hidden");
        qrContainer.classList.add("hidden");
      } else {
        // Not connected
        statusIndicator.classList.remove("bg-gray-300", "bg-green-500");
        statusIndicator.classList.add("bg-red-500");
        statusText.textContent = "Not connected to WhatsApp";
        scanQrBtn.classList.remove("hidden");
        logoutWhatsappBtn.classList.add("hidden");
      }
    } catch (error) {
      // Error checking status
      statusIndicator.classList.remove("bg-green-500", "bg-red-500");
      statusIndicator.classList.add("bg-gray-300");
      statusText.textContent = "Error checking WhatsApp connection";
      console.error("Error checking WhatsApp status:", error);
    }
  }
}

// Initialize period buttons
function initializePeriodButtons() {
  document.getElementById("last-7-days").addEventListener("click", () => {
    currentDays = 7;
    loadMessageActivity(currentDays);
  });

  document.getElementById("last-14-days").addEventListener("click", () => {
    currentDays = 14;
    loadMessageActivity(currentDays);
  });

  document.getElementById("last-30-days").addEventListener("click", () => {
    currentDays = 30;
    loadMessageActivity(currentDays);
  });
}

// Initialize date picker
function initializeDatePicker() {
  const datePicker = document.getElementById("hourly-date-picker");
  
  // Set default date (today)
  const today = new Date();
  datePicker.valueAsDate = today;
  
  // Add event listener
  datePicker.addEventListener("change", (e) => {
    loadHourlyStats(e.target.value);
  });
}

// Load all dashboard data
function loadAllDashboardData() {
  loadDashboardStats();
  loadMessageActivity(currentDays);
  loadBroadcastStatus();
  loadHourlyStats();
  loadTopRecipients();
  loadRecentBroadcasts();
}

// Set up auto-refresh
function setupAutoRefresh() {
  // Refresh dashboard data every 60 seconds
  setInterval(() => {
    loadDashboardStats();
    loadMessageActivity(currentDays);
    loadBroadcastStatus();
    loadRecentBroadcasts();
  }, 60000);
}

// API Data Loading Functions

// Load dashboard stats
async function loadDashboardStats() {
  try {
    const response = await fetch("/api/v1/dashboard/stats", {
      headers: {
        Authorization: `Bearer ${token}`,
      },
    });

    const data = await response.json();

    if (data.meta.status === "OK") {
      // Update stats cards
      document.getElementById("total-messages").textContent =
        data.data.total_messages.toLocaleString();
      document.getElementById("total-broadcasts").textContent =
        data.data.total_broadcasts.toLocaleString();
      document.getElementById("total-recipients").textContent =
        data.data.total_recipients.toLocaleString();
      document.getElementById(
        "success-rate"
      ).textContent = `${data.data.success_rate.toFixed(1)}%`;
    }
  } catch (error) {
    console.error("Error loading dashboard stats:", error);
  }
}

// Load message activity chart
async function loadMessageActivity(days = 7) {
  try {
    const response = await fetch(
      `/api/v1/dashboard/message-activity?days=${days}`,
      {
        headers: {
          Authorization: `Bearer ${token}`,
        },
      }
    );

    const data = await response.json();

    if (data.meta.status === "OK") {
      // Update chart
      const messageCtx = document
        .getElementById("message-chart")
        .getContext("2d");

      // Destroy existing chart if it exists
      if (messageChart) {
        messageChart.destroy();
      }

      // Create new chart
      messageChart = new Chart(messageCtx, {
        type: "line",
        data: {
          labels: data.data.labels,
          datasets: [
            {
              label: "Messages Sent",
              data: data.data.data,
              backgroundColor: "rgba(79, 70, 229, 0.2)",
              borderColor: "rgba(79, 70, 229, 1)",
              borderWidth: 2,
              tension: 0.3,
            },
          ],
        },
        options: {
          scales: {
            y: {
              beginAtZero: true,
              ticks: {
                precision: 0,
              },
            },
          },
          plugins: {
            title: {
              display: false,
              text: `Message Activity (Last ${days} Days)`,
            },
          },
        },
      });

      // Update period buttons
      document.querySelectorAll(".active-period").forEach((btn) => {
        btn.classList.remove(
          "active-period",
          "bg-indigo-600",
          "text-white"
        );
        btn.classList.add("bg-gray-200", "text-gray-700");
      });

      const activeBtn = document.getElementById(`last-${days}-days`);
      if (activeBtn) {
        activeBtn.classList.remove("bg-gray-200", "text-gray-700");
        activeBtn.classList.add(
          "active-period",
          "bg-indigo-600",
          "text-white"
        );
      }
    }
  } catch (error) {
    console.error("Error loading message activity:", error);
  }
}

// Load broadcast status chart
async function loadBroadcastStatus() {
  try {
    const response = await fetch("/api/v1/dashboard/broadcast-status", {
      headers: {
        Authorization: `Bearer ${token}`,
      },
    });

    const data = await response.json();

    if (data.meta.status === "OK") {
      // Update chart
      const broadcastCtx = document
        .getElementById("broadcast-chart")
        .getContext("2d");

      // Destroy existing chart if it exists
      if (broadcastChart) {
        broadcastChart.destroy();
      }

      // Create new chart
      broadcastChart = new Chart(broadcastCtx, {
        type: "doughnut",
        data: {
          labels: data.data.labels,
          datasets: [
            {
              data: data.data.data,
              backgroundColor: [
                "rgba(16, 185, 129, 0.7)", // Success - Green
                "rgba(245, 158, 11, 0.7)", // Pending - Yellow
                "rgba(239, 68, 68, 0.7)", // Failed - Red
              ],
              borderColor: [
                "rgba(16, 185, 129, 1)",
                "rgba(245, 158, 11, 1)",
                "rgba(239, 68, 68, 1)",
              ],
              borderWidth: 1,
            },
          ],
        },
        options: {
          responsive: true,
          plugins: {
            legend: {
              position: "bottom",
            },
          },
        },
      });
    }
  } catch (error) {
    console.error("Error loading broadcast status:", error);
  }
}

// Load hourly message stats
async function loadHourlyStats(dateStr = "") {
  try {
    const url = dateStr
      ? `/api/v1/dashboard/hourly-stats?date=${dateStr}`
      : "/api/v1/dashboard/hourly-stats";

    const response = await fetch(url, {
      headers: {
        Authorization: `Bearer ${token}`,
      },
    });

    const data = await response.json();

    if (data.meta.status === "OK") {
      // Update chart
      const hourlyCtx = document
        .getElementById("hourly-chart")
        .getContext("2d");

      // Destroy existing chart if it exists
      if (hourlyChart) {
        hourlyChart.destroy();
      }

      // Create new chart
      hourlyChart = new Chart(hourlyCtx, {
        type: "bar",
        data: {
          labels: data.data.labels,
          datasets: [
            {
              label: "Messages",
              data: data.data.data,
              backgroundColor: "rgba(59, 130, 246, 0.5)",
              borderColor: "rgba(59, 130, 246, 1)",
              borderWidth: 1,
            },
          ],
        },
        options: {
          scales: {
            y: {
              beginAtZero: true,
              ticks: {
                precision: 0,
              },
            },
          },
          plugins: {
            title: {
              display: true,
              text: `Hourly Distribution (${data.data.date})`,
            },
          },
        },
      });

      // Update date picker value
      document.getElementById("hourly-date-picker").value =
        data.data.date;
    }
  } catch (error) {
    console.error("Error loading hourly stats:", error);
  }
}

// Load top recipients
async function loadTopRecipients() {
  try {
    const response = await fetch("/api/v1/dashboard/top-recipients", {
      headers: {
        Authorization: `Bearer ${token}`,
      },
    });

    const data = await response.json();

    if (data.meta.status === "OK") {
      const tableBody = document.getElementById("top-recipients-table");

      if (data.data.length === 0) {
        tableBody.innerHTML = `<tr><td colspan="2" class="px-6 py-4 text-center text-sm text-gray-500">No data available</td></tr>`;
        return;
      }

      // Generate table rows
      let html = "";
      data.data.forEach((recipient) => {
        html += `
          <tr>
            <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">${recipient.recipient_number}</td>
            <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">${recipient.message_count}</td>
          </tr>
        `;
      });

      tableBody.innerHTML = html;
    }
  } catch (error) {
    console.error("Error loading top recipients:", error);
  }
}

// Load recent broadcasts
async function loadRecentBroadcasts() {
  try {
    const response = await fetch(
      "/api/v1/dashboard/recent-broadcasts",
      {
        headers: {
          Authorization: `Bearer ${token}`,
        },
      }
    );

    const data = await response.json();

    if (data.meta.status === "OK") {
      const tableBody = document.getElementById(
        "recent-broadcasts-table"
      );

      if (data.data.length === 0) {
        tableBody.innerHTML = `<tr><td colspan="6" class="px-6 py-4 text-center text-sm text-gray-500">No broadcasts available</td></tr>`;
        return;
      }

      // Generate table rows
      let html = "";
      data.data.forEach((broadcast) => {
        // Format date
        const date = new Date(broadcast.created_at);
        const formattedDate = date.toLocaleDateString("en-US", {
          year: "numeric",
          month: "short",
          day: "numeric",
          hour: "2-digit",
          minute: "2-digit",
        });

        // Status badge color
        let statusClass = "bg-gray-100 text-gray-800";
        if (
          broadcast.broadcast_status === "Success" ||
          broadcast.broadcast_status === "Completed"
        ) {
          statusClass = "bg-green-100 text-green-800";
        } else if (broadcast.broadcast_status === "Failed") {
          statusClass = "bg-red-100 text-red-800";
        } else if (
          broadcast.broadcast_status === "Processing" ||
          broadcast.broadcast_status === "Starting"
        ) {
          statusClass = "bg-yellow-100 text-yellow-800";
        }

        html += `
          <tr>
            <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">${
              broadcast.client_info
            }</td>
            <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">${formattedDate}</td>
            <td class="px-6 py-4 whitespace-nowrap text-sm">
              <span class="px-2 inline-flex text-xs leading-5 font-semibold rounded-full ${statusClass}">
                ${broadcast.broadcast_status || "Pending"}
              </span>
            </td>
            <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">${
              broadcast.total_recipients
            }</td>
            <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">${broadcast.success_rate.toFixed(
              1
            )}%</td>
            <td class="px-6 py-4 whitespace-nowrap text-sm">
              <a href="/admin/broadcast/recipients?id=${
                broadcast.id
              }" class="text-indigo-600 hover:text-indigo-900">
                <i class="fas fa-eye mr-1"></i> View
              </a>
            </td>
          </tr>
        `;
      });

      tableBody.innerHTML = html;
    }
  } catch (error) {
    console.error("Error loading recent broadcasts:", error);
  }
}
