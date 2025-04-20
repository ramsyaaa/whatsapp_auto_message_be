// Reports JavaScript Logic

// Global variables
let token;
let reportTable;
let selectedBroadcast;

// Initialize the reports page when the DOM is loaded
document.addEventListener("DOMContentLoaded", function () {
  initializeReportsPage();
});

// Main initialization function
function initializeReportsPage() {
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

  // Load broadcasts for the dropdown
  loadBroadcasts();

  // Initialize report generation
  initializeReportGeneration();
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

// Load broadcasts for the dropdown
async function loadBroadcasts() {
  try {
    const response = await fetch("/api/v1/broadcast", {
      headers: {
        Authorization: `Bearer ${token}`,
      },
    });

    const data = await response.json();

    if (data.meta.status === "OK") {
      const broadcastSelect = document.getElementById("broadcast-select");
      
      if (data.data.length === 0) {
        broadcastSelect.innerHTML = '<option value="">No broadcasts available</option>';
        return;
      }

      // Sort broadcasts by date (newest first)
      const sortedBroadcasts = data.data.sort((a, b) => {
        return new Date(b.created_at) - new Date(a.created_at);
      });

      // Generate options
      let options = '<option value="">Select a broadcast</option>';
      sortedBroadcasts.forEach((broadcast) => {
        const date = new Date(broadcast.created_at);
        const formattedDate = date.toLocaleDateString("en-US", {
          year: "numeric",
          month: "short",
          day: "numeric",
        });
        
        options += `<option value="${broadcast.id}">${broadcast.client_info} - ${formattedDate}</option>`;
      });

      broadcastSelect.innerHTML = options;
    } else {
      document.getElementById("broadcast-select").innerHTML = 
        '<option value="">Error loading broadcasts</option>';
    }
  } catch (error) {
    console.error("Error loading broadcasts:", error);
    document.getElementById("broadcast-select").innerHTML = 
      '<option value="">Error loading broadcasts</option>';
  }
}

// Initialize report generation
function initializeReportGeneration() {
  const generateReportBtn = document.getElementById("generate-report-btn");
  const broadcastSelect = document.getElementById("broadcast-select");

  generateReportBtn.addEventListener("click", () => {
    const broadcastId = broadcastSelect.value;
    
    if (!broadcastId) {
      alert("Please select a broadcast");
      return;
    }

    selectedBroadcast = broadcastId;
    
    // Get selected report type
    const reportType = document.querySelector('input[name="report-type"]:checked').id;
    
    // Generate report based on type
    switch (reportType) {
      case "report-type-summary":
        generateSummaryReport(broadcastId);
        break;
      case "report-type-detailed":
        generateDetailedReport(broadcastId, "all");
        break;
      case "report-type-success":
        generateDetailedReport(broadcastId, "success");
        break;
      case "report-type-failed":
        generateDetailedReport(broadcastId, "failed");
        break;
      default:
        generateSummaryReport(broadcastId);
    }
  });
}

// Generate summary report
async function generateSummaryReport(broadcastId) {
  try {
    // Show loading state
    document.getElementById("report-results").classList.remove("hidden");
    document.getElementById("summary-section").classList.remove("hidden");
    
    // Reset summary values
    document.getElementById("summary-total").textContent = "--";
    document.getElementById("summary-success-rate").textContent = "--";
    document.getElementById("summary-date").textContent = "--";
    document.getElementById("summary-success").textContent = "--";
    document.getElementById("summary-failed").textContent = "--";
    document.getElementById("summary-pending").textContent = "--";
    
    // Fetch broadcast details
    const response = await fetch(`/api/v1/broadcast/${broadcastId}`, {
      headers: {
        Authorization: `Bearer ${token}`,
      },
    });

    const data = await response.json();

    if (data.meta.status === "OK") {
      const broadcast = data.data;
      
      // Set broadcast info
      document.getElementById("broadcast-info").textContent = 
        `${broadcast.client_info} - ${broadcast.broadcast_type === "pecatu" ? "Pecatu" : "Regular"}`;
      
      // Format date
      const date = new Date(broadcast.created_at);
      const formattedDate = date.toLocaleDateString("en-US", {
        year: "numeric",
        month: "short",
        day: "numeric",
      });
      
      // Update summary section
      document.getElementById("summary-total").textContent = broadcast.total_recipients || 0;
      document.getElementById("summary-success-rate").textContent = `${(broadcast.success_rate || 0).toFixed(1)}%`;
      document.getElementById("summary-date").textContent = formattedDate;
      
      // Fetch recipients to get counts by status
      const recipientsResponse = await fetch(`/api/v1/broadcast/recipients/${broadcastId}`, {
        headers: {
          Authorization: `Bearer ${token}`,
        },
      });

      const recipientsData = await recipientsResponse.json();
      
      if (recipientsData.meta.status === "OK") {
        const recipients = recipientsData.data.recipients || [];
        
        // Count by status
        const successCount = recipients.filter(r => r.broadcast_status === "Success").length;
        const failedCount = recipients.filter(r => r.broadcast_status === "Failed").length;
        const pendingCount = recipients.filter(r => r.broadcast_status === "Pending").length;
        
        document.getElementById("summary-success").textContent = successCount;
        document.getElementById("summary-failed").textContent = failedCount;
        document.getElementById("summary-pending").textContent = pendingCount;
      }
      
      // Generate detailed report table
      generateDetailedReport(broadcastId, "all", false);
    }
  } catch (error) {
    console.error("Error generating summary report:", error);
    alert("Error generating report. Please try again.");
  }
}

// Generate detailed report
async function generateDetailedReport(broadcastId, filter = "all", showSummary = true) {
  try {
    // Show loading state
    document.getElementById("report-results").classList.remove("hidden");
    
    if (showSummary) {
      document.getElementById("summary-section").classList.add("hidden");
    }
    
    // Fetch broadcast details for the header
    const broadcastResponse = await fetch(`/api/v1/broadcast/${broadcastId}`, {
      headers: {
        Authorization: `Bearer ${token}`,
      },
    });

    const broadcastData = await broadcastResponse.json();
    
    if (broadcastData.meta.status === "OK") {
      const broadcast = broadcastData.data;
      
      // Set broadcast info
      document.getElementById("broadcast-info").textContent = 
        `${broadcast.client_info} - ${broadcast.broadcast_type === "pecatu" ? "Pecatu" : "Regular"}`;
    }
    
    // Fetch recipients
    const response = await fetch(`/api/v1/broadcast/recipients/${broadcastId}`, {
      headers: {
        Authorization: `Bearer ${token}`,
      },
    });

    const data = await response.json();

    if (data.meta.status === "OK") {
      let recipients = data.data.recipients || [];
      
      // Filter recipients based on the filter parameter
      if (filter === "success") {
        recipients = recipients.filter(r => r.broadcast_status === "Success");
      } else if (filter === "failed") {
        recipients = recipients.filter(r => r.broadcast_status === "Failed");
      }
      
      // Prepare data for DataTable
      const tableData = recipients.map((recipient, index) => {
        // Status badge color
        let statusClass = "bg-gray-100 text-gray-800";
        if (recipient.broadcast_status === "Success") {
          statusClass = "bg-green-100 text-green-800";
        } else if (recipient.broadcast_status === "Pending") {
          statusClass = "bg-yellow-100 text-yellow-800";
        } else if (recipient.broadcast_status === "Failed") {
          statusClass = "bg-red-100 text-red-800";
        }
        
        // Format broadcast_at date if available
        let broadcastAtText = "-";
        if (recipient.broadcast_at) {
          const date = new Date(recipient.broadcast_at);
          broadcastAtText = date.toLocaleString();
        }
        
        // Get message content
        const messageContent = broadcastData.data.broadcast_message || "-";
        
        return [
          index + 1, // Row number
          recipient.whatsapp_number || "-",
          recipient.recipient_name || "-",
          `<span class="px-2 inline-flex text-xs leading-5 font-semibold rounded-full ${statusClass}">${recipient.broadcast_status || "Pending"}</span>`,
          broadcastAtText,
          messageContent
        ];
      });
      
      // Initialize or refresh DataTable
      if (reportTable) {
        reportTable.destroy();
      }
      
      // Get report type for title
      let reportTitle = "Detailed Report";
      if (filter === "success") {
        reportTitle = "Success Report";
      } else if (filter === "failed") {
        reportTitle = "Failed Report";
      }
      
      // Initialize DataTable with export buttons
      reportTable = $("#report-table").DataTable({
        data: tableData,
        dom: 'Bfrtip',
        buttons: [
          {
            extend: 'excel',
            text: '<i class="fas fa-file-excel mr-1"></i> Export to Excel',
            title: `${broadcastData.data.client_info} - ${reportTitle} - ${new Date().toLocaleDateString()}`,
            className: 'export-excel-button'
          },
          {
            extend: 'csv',
            text: '<i class="fas fa-file-csv mr-1"></i> Export to CSV',
            title: `${broadcastData.data.client_info} - ${reportTitle} - ${new Date().toLocaleDateString()}`,
            className: 'export-csv-button'
          },
          {
            extend: 'print',
            text: '<i class="fas fa-print mr-1"></i> Print',
            title: `${broadcastData.data.client_info} - ${reportTitle} - ${new Date().toLocaleDateString()}`,
            className: 'print-button'
          }
        ],
        pageLength: 25,
        language: {
          search: "Search:",
          lengthMenu: "Show _MENU_ entries",
          info: "Showing _START_ to _END_ of _TOTAL_ entries",
          paginate: {
            first: "First",
            last: "Last",
            next: "Next",
            previous: "Previous",
          },
          emptyTable: "No data available"
        }
      });
    }
  } catch (error) {
    console.error("Error generating detailed report:", error);
    alert("Error generating report. Please try again.");
  }
}
