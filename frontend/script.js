const API_BASE_URL = 'http://localhost:8080/api';

let encodedImageBlob = null;

// Tab switching
function switchTab(tabName) {
    // Hide all tabs
    document.querySelectorAll('.tab-content').forEach(tab => {
        tab.classList.remove('active');
    });
    
    // Remove active from all buttons
    document.querySelectorAll('.tab-button').forEach(btn => {
        btn.classList.remove('active');
    });
    
    // Show selected tab
    document.getElementById(`${tabName}-tab`).classList.add('active');
    event.target.classList.add('active');
    
    // Reset forms and results
    document.getElementById('encode-result').style.display = 'none';
    document.getElementById('decode-result').style.display = 'none';
}

// Image preview functionality
function setupImagePreview(inputId, previewId) {
    const input = document.getElementById(inputId);
    const preview = document.getElementById(previewId);
    
    input.addEventListener('change', function(e) {
        const file = e.target.files[0];
        if (file) {
            const reader = new FileReader();
            reader.onload = function(e) {
                preview.innerHTML = `<img src="${e.target.result}" alt="Preview">`;
            };
            reader.readAsDataURL(file);
        } else {
            preview.innerHTML = '';
        }
    });
}

// Character counter
function setupCharCounter() {
    const textarea = document.getElementById('message');
    const counter = document.getElementById('char-count');
    
    textarea.addEventListener('input', function() {
        counter.textContent = textarea.value.length;
    });
}

// Show notification
function showNotification(message, type = 'success') {
    const notification = document.getElementById('notification');
    notification.textContent = message;
    notification.className = `notification show ${type}`;
    
    setTimeout(() => {
        notification.classList.remove('show');
    }, 3000);
}

// Loading state for buttons
function setButtonLoading(button, isLoading) {
    const btnText = button.querySelector('.btn-text');
    const loader = button.querySelector('.loader');
    
    if (isLoading) {
        btnText.style.display = 'none';
        loader.style.display = 'block';
        button.disabled = true;
    } else {
        btnText.style.display = 'inline';
        loader.style.display = 'none';
        button.disabled = false;
    }
}

// Handle encode form submission
document.getElementById('encode-form').addEventListener('submit', async function(e) {
    e.preventDefault();
    
    const imageInput = document.getElementById('encode-image');
    const messageInput = document.getElementById('message');
    const submitBtn = this.querySelector('button[type="submit"]');
    const resultDiv = document.getElementById('encode-result');
    
    if (!imageInput.files[0]) {
        showNotification('Please select an image', 'error');
        return;
    }
    
    if (!messageInput.value.trim()) {
        showNotification('Please enter a message', 'error');
        return;
    }
    
    setButtonLoading(submitBtn, true);
    resultDiv.style.display = 'none';
    
    try {
        const formData = new FormData();
        formData.append('image', imageInput.files[0]);
        formData.append('message', messageInput.value);
        
        // Add password if provided
        const passwordInput = document.getElementById('encode-password');
        if (passwordInput.value) {
            formData.append('password', passwordInput.value);
        }
        
        const response = await fetch(`${API_BASE_URL}/encode`, {
            method: 'POST',
            body: formData
        });
        
        if (!response.ok) {
            const errorData = await response.json();
            throw new Error(errorData.error || 'Failed to encode message');
        }
        
        // Get the encoded image as blob
        encodedImageBlob = await response.blob();
        
        // Show success message
        resultDiv.style.display = 'block';
        showNotification('Message encoded successfully!', 'success');
        
    } catch (error) {
        console.error('Encode error:', error);
        showNotification(error.message || 'Failed to encode message', 'error');
    } finally {
        setButtonLoading(submitBtn, false);
    }
});

// Handle download button
document.getElementById('download-btn').addEventListener('click', function() {
    if (!encodedImageBlob) {
        showNotification('No encoded image available', 'error');
        return;
    }
    
    const url = URL.createObjectURL(encodedImageBlob);
    const a = document.createElement('a');
    a.href = url;
    a.download = 'encoded_image.png';
    document.body.appendChild(a);
    a.click();
    document.body.removeChild(a);
    URL.revokeObjectURL(url);
    
    showNotification('Image downloaded!', 'success');
});

// Handle decode form submission
document.getElementById('decode-form').addEventListener('submit', async function(e) {
    e.preventDefault();
    
    const imageInput = document.getElementById('decode-image');
    const submitBtn = this.querySelector('button[type="submit"]');
    const resultDiv = document.getElementById('decode-result');
    const messageOutput = document.getElementById('decoded-message');
    
    if (!imageInput.files[0]) {
        showNotification('Please select an image', 'error');
        return;
    }
    
    setButtonLoading(submitBtn, true);
    resultDiv.style.display = 'none';
    
    try {
        const formData = new FormData();
        formData.append('image', imageInput.files[0]);
        
        // Add password if provided
        const passwordInput = document.getElementById('decode-password');
        if (passwordInput.value) {
            formData.append('password', passwordInput.value);
        }
        
        const response = await fetch(`${API_BASE_URL}/decode`, {
            method: 'POST',
            body: formData
        });
        
        if (!response.ok) {
            const errorData = await response.json();
            throw new Error(errorData.error || 'Failed to decode message');
        }
        
        const data = await response.json();
        
        if (!data.message || data.message.trim() === '') {
            showNotification('No hidden message found in this image', 'error');
            return;
        }
        
        // Display the decoded message
        messageOutput.textContent = data.message;
        resultDiv.style.display = 'block';
        showNotification('Message decoded successfully!', 'success');
        
    } catch (error) {
        console.error('Decode error:', error);
        showNotification(error.message || 'Failed to decode message', 'error');
    } finally {
        setButtonLoading(submitBtn, false);
    }
});

// Handle copy button
document.getElementById('copy-btn').addEventListener('click', function() {
    const messageOutput = document.getElementById('decoded-message');
    const text = messageOutput.textContent;
    
    navigator.clipboard.writeText(text).then(() => {
        showNotification('Message copied to clipboard!', 'success');
    }).catch(err => {
        console.error('Failed to copy:', err);
        showNotification('Failed to copy message', 'error');
    });
});

// Initialize on page load
document.addEventListener('DOMContentLoaded', function() {
    setupImagePreview('encode-image', 'encode-preview');
    setupImagePreview('decode-image', 'decode-preview');
    setupCharCounter();
});
