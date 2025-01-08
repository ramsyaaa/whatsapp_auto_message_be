document.addEventListener("DOMContentLoaded", function () {
  const logFileSelect = document.getElementById("logFileSelect");
  const methodFilter = document.getElementById("methodFilter");
  const statusFilter = document.getElementById("statusFilter");
  const pathFilter = document.getElementById("pathFilter");
  const logList = document.getElementById("logList");
  const exportButton = document.getElementById("exportButton");
  const prevPageButton = document.getElementById("prevPage");
  const nextPageButton = document.getElementById("nextPage");
  const pageNumberElement = document.getElementById("pageNumber");

  let logs = []; // Store logs fetched from the server
  let currentLogFile = ""; // Store the selected log file
  let currentPage = 1;
  const rowsPerPage = 10;
  let totalLogs = 0;

  // Fetch the available log files
  function fetchLogFiles() {
    fetch("/logs")
      .then((response) => response.json())
      .then((data) => {
        const files = data.files;
        files.forEach((file) => {
          const option = document.createElement("option");
          option.value = file;
          option.textContent = file;
          logFileSelect.appendChild(option);
        });

        // Initialize Select2 on the log file select element after options are loaded
        $(logFileSelect).select2(); // Initialize Select2
      })
      .catch((error) => console.error("Error fetching log files:", error));
  }

  // Function to fetch logs from the selected log file, including pagination
  function fetchLogs() {
    if (!currentLogFile) {
      return; // Do nothing if no log file is selected
    }

    fetch(`/logs/${currentLogFile}?page=${currentPage}&limit=${rowsPerPage}`)
      .then((response) => response.json())
      .then((data) => {
        logs = data.logs; // Store logs
        totalLogs = data.total;
        loadLogs(); // Load the logs into the UI
        updatePagination();
        exportButton.classList.remove("hidden"); // Show the export button
      })
      .catch((error) => console.error("Error fetching logs:", error));
  }

  // Function to load logs into the table UI with filtering applied
  function loadLogs() {
    logList.innerHTML = ""; // Clear current list

    const filteredLogs = logs.filter((log) => {
      const methodMatch = methodFilter.value
        ? log.method === methodFilter.value
        : true;
      const statusMatch = statusFilter.value
        ? log.status_code === statusFilter.value
        : true;
      const pathMatch = pathFilter.value
        ? log.path.includes(pathFilter.value)
        : true;

      return methodMatch && statusMatch && pathMatch;
    });

    filteredLogs.forEach((log) => {
      const row = document.createElement("tr");
      row.classList.add("border-t", "hover:bg-gray-100");

      // Apply color based on status code
      let statusClass = "bg-gray-100"; // Default
      if (log.status_code === "200") {
        statusClass = "bg-green-200 text-green-800";
      } else if (["400", "404"].includes(log.status_code)) {
        statusClass = "bg-yellow-200 text-yellow-800";
      } else if (log.status_code === "500") {
        statusClass = "bg-red-200 text-red-800";
      }

      row.innerHTML = `
        <td class="px-6 py-4 text-sm text-gray-600">${log.timestamp}</td>
        <td class="px-6 py-4 text-sm ${statusClass}">${log.status_code}</td>
        <td class="px-6 py-4 text-sm text-gray-600">${log.latency}</td>
        <td class="px-6 py-4 text-sm text-gray-600">${log.ip}</td>
        <td class="px-6 py-4 text-sm text-gray-600">${log.method}</td>
        <td class="px-6 py-4 text-sm text-gray-600">${log.path}</td>
        <td class="px-6 py-4 text-sm text-gray-600">${log.additional_info}</td>
      `;
      logList.appendChild(row);
    });
  }

  // Function to update pagination controls
  function updatePagination() {
    pageNumberElement.textContent = `Page ${currentPage}`;
    prevPageButton.disabled = currentPage === 1;
    nextPageButton.disabled = currentPage * rowsPerPage >= totalLogs;
  }

  // Handle pagination events
  prevPageButton.addEventListener("click", function () {
    if (currentPage > 1) {
      currentPage--;
      fetchLogs();
    }
  });

  nextPageButton.addEventListener("click", function () {
    if (currentPage * rowsPerPage < totalLogs) {
      currentPage++;
      fetchLogs();
    }
  });

  // Handle file selection with Select2
  $(logFileSelect).on("change", function () {
    currentLogFile = $(logFileSelect).val(); // Get selected value
    fetchLogs(); // Fetch logs when a file is selected
  });

  // Event listeners for filters
  methodFilter.addEventListener("change", loadLogs);
  statusFilter.addEventListener("change", loadLogs);
  pathFilter.addEventListener("input", loadLogs);

  // Function to export data to Excel
  function exportToExcel() {
    const wb = XLSX.utils.book_new();
    const ws = XLSX.utils.table_to_sheet(document.getElementById("logTable"));
    XLSX.utils.book_append_sheet(wb, ws, "Logs");

    const logFileName = currentLogFile || "logs";
    const excelFilename = `${logFileName.replace(".log", "")}.xlsx`;
    XLSX.writeFile(wb, excelFilename);
  }

  // Add event listener for export button
  exportButton.addEventListener("click", exportToExcel);

  // Initial fetch of log files and loading the first file
  fetchLogFiles();
});
