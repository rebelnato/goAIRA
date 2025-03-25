function showMessage(id, event) {
    event.preventDefault(); // Prevent page refresh

    // Hide all .message divs
    document.querySelectorAll('.message').forEach(div => {
        div.style.display = "none";
    });

    // Show the selected message div
    let selectedDiv = document.getElementById(id);
    if (selectedDiv) {
        selectedDiv.style.display = "flex";
    }

    // Remove 'selected' class from all icons
    document.querySelectorAll('.vertical-nav ul li a').forEach(link => {
        link.classList.remove('selected');
    });

    // Add 'selected' class to the clicked icon
    event.currentTarget.classList.add('selected');
}

document.addEventListener("DOMContentLoaded", function () {
    document.querySelectorAll(".submit1").forEach(button => {
        button.addEventListener("click", function (event) {
            event.preventDefault(); // Prevent default form submission

            const parentDiv = this.closest(".message"); // Get the closest message div
            const form = parentDiv.querySelector("form"); // Find the form
            const outputDiv = parentDiv.querySelector(".output"); // Find the output div
            outputDiv.innerHTML = "<p>Loading...</p>"; // Show loading state
            
            // Hide form & show output div
            form.style.display = "none";
            outputDiv.style.display = "flex"; // Ensure it's visible

            // Define API URL & method based on selected message div
            let apiUrl = "";
            let method = "";

            switch (parentDiv.id) {
                case "health":
                    apiUrl = "http://localhost:8080/health";
                    method = "GET";
                    break;
                case "create":
                    apiUrl = "http://localhost:8080/createincident";
                    method = "POST";
                    break;
                case "get":
                    apiUrl = "http://localhost:8080/getincident";
                    method = "GET";
                    break;
                case "update":
                    apiUrl = "http://localhost:8080/updateincident";
                    method = "PATCH";
                    break;
                default:
                    outputDiv.innerHTML = "<span style='color: red;'>Invalid action</span>";
                    return;
            }

            // Collect headers from form inputs inside the selected `.message` div
            let headersData = {};
            parentDiv.querySelectorAll(".header-input").forEach(input => {
                if (input.value.trim() !== "") {
                    headersData[input.name] = input.value; // Add input name & value as header
                }
            });

            // Make the API request
            fetch(apiUrl, {
                method: method,
                headers: {
                    "Content-Type": "application/json",
                    ...headersData
                }
            })
            .then(response => response.json()) // Convert response to JSON
            .then(data => {
                outputDiv.innerHTML = `
                    <div class="response-container">
                        <pre>${JSON.stringify(data, null, 2)}</pre>
                        <button class="back-button">New</button>
                    </div>
                `; // Wrap response in `.output`
                
                // Add event listener to "Back" button
                parentDiv.querySelector(".back-button").addEventListener("click", function () {
                    outputDiv.style.display = "none"; // Hide output
                    form.style.display = "block"; // Show form again
                });
            })
            .catch(error => {
                outputDiv.innerHTML = `
                    <div class="response-container">
                        <span style="color: red;">Error: ${error.message}</span>
                        <button class="back-button">Back</button>
                    </div>
                `;

                // Add event listener to "Back" button
                parentDiv.querySelector(".back-button").addEventListener("click", function () {
                    outputDiv.style.display = "none"; // Hide output
                    form.style.display = "block"; // Show form again
                });
            });
        });
    });
});
