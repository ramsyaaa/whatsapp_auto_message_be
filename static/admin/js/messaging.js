// Messaging JavaScript Logic

// Global variables
let token;

// Initialize the messaging page when the DOM is loaded
document.addEventListener("DOMContentLoaded", function () {
  initializeMessagingPage();
});

// Main initialization function
function initializeMessagingPage() {
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

  // Initialize message form
  initializeMessageForm();

  // Load recent messages
  loadRecentMessages();
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

  checkWhatsAppStatus();

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
      } else {
        // Not connected
        statusIndicator.classList.remove("bg-gray-300", "bg-green-500");
        statusIndicator.classList.add("bg-red-500");
        statusText.textContent =
          "Not connected to WhatsApp. Please connect from the Dashboard.";
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

// Initialize message form
function initializeMessageForm() {
  const messageForm = document.getElementById("message-form");
  const errorMessage = document.getElementById("error-message");
  const successMessage = document.getElementById("success-message");
  const sendText = document.getElementById("send-text");
  const sendSpinner = document.getElementById("send-spinner");

  messageForm.addEventListener("submit", async function (e) {
    e.preventDefault();

    // Show loading state
    sendText.textContent = "Sending...";
    sendSpinner.classList.remove("hidden");
    errorMessage.classList.add("hidden");
    successMessage.classList.add("hidden");

    const phoneNumber = document.getElementById("phone-number").value;
    const message = document.getElementById("message").value;

    try {
      const response = await fetch("/api/v1/send-message", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${token}`,
        },
        body: JSON.stringify({
          number: phoneNumber,
          message: message,
        }),
      });

      const data = await response.json();

      if (response.ok) {
        // Show success message
        successMessage.textContent = "Message sent successfully!";
        successMessage.classList.remove("hidden");

        // Reset form
        messageForm.reset();

        // Add to recent messages
        addRecentMessage(
          phoneNumber,
          message,
          "Sent",
          new Date().toLocaleString()
        );
      } else {
        // Show error message
        errorMessage.textContent =
          data.error || "Failed to send message. Please try again.";
        errorMessage.classList.remove("hidden");
      }
    } catch (error) {
      // Show error message
      errorMessage.textContent = "An error occurred. Please try again.";
      errorMessage.classList.remove("hidden");
      console.error("Error sending message:", error);
    } finally {
      // Reset send button
      sendText.textContent = "Send Message";
      sendSpinner.classList.add("hidden");
    }
  });
}

// Function to load recent messages from the database
async function loadRecentMessages() {
  const recentMessages = document.getElementById("recent-messages");
  recentMessages.innerHTML =
    '<tr><td colspan="4" class="px-6 py-4 text-center text-sm text-gray-500">Loading messages...</td></tr>';

  try {
    const response = await fetch("/api/v1/recent-messages", {
      headers: {
        Authorization: `Bearer ${token}`,
      },
    });

    const data = await response.json();

    if (response.ok && data.data && data.data.length > 0) {
      recentMessages.innerHTML = "";

      data.data.forEach((message) => {
        // Create new row
        const row = document.createElement("tr");

        // Truncate message if it's too long
        const truncatedMessage =
          message.message_content.length > 50
            ? message.message_content.substring(0, 50) + "..."
            : message.message_content;

        // Format date
        const sentAt = new Date(message.sent_at).toLocaleString();

        // Status badge color
        let statusClass = "bg-green-100 text-green-800";
        if (message.status === "Failed") {
          statusClass = "bg-red-100 text-red-800";
        } else if (message.status === "Pending") {
          statusClass = "bg-yellow-100 text-yellow-800";
        }

        // Set row content
        row.innerHTML = `
          <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
              ${message.recipient_number}
          </td>
          <td class="px-6 py-4 text-sm text-gray-500">
              ${truncatedMessage}
          </td>
          <td class="px-6 py-4 whitespace-nowrap text-sm">
              <span class="px-2 inline-flex text-xs leading-5 font-semibold rounded-full ${statusClass}">
                  ${message.status}
              </span>
          </td>
          <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
              ${sentAt}
          </td>
        `;

        // Add row to table
        recentMessages.appendChild(row);
      });
    } else {
      recentMessages.innerHTML =
        '<tr><td colspan="4" class="px-6 py-4 text-center text-sm text-gray-500">No recent messages</td></tr>';
    }
  } catch (error) {
    console.error("Error loading recent messages:", error);
    recentMessages.innerHTML =
      '<tr><td colspan="4" class="px-6 py-4 text-center text-sm text-gray-500">Error loading messages</td></tr>';
  }
}

// Function to add a message to the recent messages table (for immediate display after sending)
function addRecentMessage(recipient, message, status, time) {
  const recentMessages = document.getElementById("recent-messages");

  // Remove "No recent messages" row if it exists
  if (recentMessages.querySelector('td[colspan="4"]')) {
    recentMessages.innerHTML = "";
  }

  // Create new row
  const row = document.createElement("tr");

  // Truncate message if it's too long
  const truncatedMessage =
    message.length > 50 ? message.substring(0, 50) + "..." : message;

  // Set row content
  row.innerHTML = `
    <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
        ${recipient}
    </td>
    <td class="px-6 py-4 text-sm text-gray-500">
        ${truncatedMessage}
    </td>
    <td class="px-6 py-4 whitespace-nowrap text-sm">
        <span class="px-2 inline-flex text-xs leading-5 font-semibold rounded-full bg-green-100 text-green-800">
            ${status}
        </span>
    </td>
    <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
        ${time}
    </td>
  `;

  // Add row to table
  recentMessages.insertBefore(row, recentMessages.firstChild);
}
