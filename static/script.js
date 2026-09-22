// DOM Elements
const formSection = document.getElementById('formSection');
const resultSection = document.getElementById('resultSection');
const loveForm = document.getElementById('loveForm');
const name1Input = document.getElementById('name1');
const name2Input = document.getElementById('name2');
const submitBtn = document.getElementById('submitBtn');
const resetBtn = document.getElementById('resetBtn');
const errorMessage = document.getElementById('errorMessage');
const loadingOverlay = document.getElementById('loadingOverlay');
const loadingQuote = document.getElementById('loadingQuote');
const percentageNumber = document.getElementById('percentageNumber');
const resultMessage = document.getElementById('resultMessage');

// Cheesy loading quotes
const loadingQuotes = [
    "💘 Consulting the love gods...",
    "💌 Asking Cupid for permission...",
    "🔮 Calculating the cosmic chemistry...",
    "💕 Checking the astrological charts...",
    "✨ Aligning your romantic chakras...",
    "🦋 Measuring heart palpitations...",
    "🎀 Analyzing romantic wavelengths..."
];

let quoteInterval = null;

// Form submission handler
loveForm.addEventListener('submit', async (e) => {
    e.preventDefault();
    await handleCalculate();
});

// Reset button handler
resetBtn.addEventListener('click', handleReset);

/**
 * Validate input on the frontend.
 */
function validateInput(name1, name2) {
    const trimmedName1 = name1.trim();
    const trimmedName2 = name2.trim();

    if (!trimmedName1) {
        return "Oops! Please enter your name. ❤️";
    }
    if (!trimmedName2) {
        return "Oops! Please enter their name. 💘";
    }
    if (trimmedName1.length > 50) {
        return "First name is too long (max 50 characters). 😭";
    }
    if (trimmedName2.length > 50) {
        return "Second name is too long (max 50 characters). 😭";
    }

    return null;
}

/**
 * Show error message on the form.
 */
function showError(message) {
    errorMessage.textContent = message;
    errorMessage.style.display = 'block';
}

/**
 * Clear error message.
 */
function clearError() {
    errorMessage.textContent = '';
    errorMessage.style.display = 'none';
}

/**
 * Show loading overlay with rotating cheesy quotes.
 */
function showLoading() {
    let quoteIndex = 0;
    loadingQuote.textContent = loadingQuotes[0];
    loadingOverlay.style.display = 'flex';

    quoteInterval = setInterval(() => {
        quoteIndex = (quoteIndex + 1) % loadingQuotes.length;
        loadingQuote.textContent = loadingQuotes[quoteIndex];
    }, 800);
}

/**
 * Hide loading overlay.
 */
function hideLoading() {
    if (quoteInterval) {
        clearInterval(quoteInterval);
        quoteInterval = null;
    }
    loadingOverlay.style.display = 'none';
}

/**
 * Animate percentage counter from 0 to the target value.
 */
function animatePercentage(targetPercentage, duration = 1200) {
    return new Promise((resolve) => {
        const startTime = Date.now();
        const startValue = 0;

        function update() {
            const elapsed = Date.now() - startTime;
            const progress = Math.min(elapsed / duration, 1);

            // Easing: ease-out cubic
            const easeProgress = 1 - Math.pow(1 - progress, 3);
            const currentValue = Math.floor(startValue + (targetPercentage - startValue) * easeProgress);

            percentageNumber.textContent = currentValue;

            if (progress < 1) {
                requestAnimationFrame(update);
            } else {
                percentageNumber.textContent = targetPercentage;
                resolve();
            }
        }

        requestAnimationFrame(update);
    });
}

/**
 * Handle the calculate button click.
 */
async function handleCalculate() {
    clearError();

    const name1 = name1Input.value;
    const name2 = name2Input.value;

    // Frontend validation
    const validationError = validateInput(name1, name2);
    if (validationError) {
        showError(validationError);
        return;
    }

    // Disable submit button
    submitBtn.disabled = true;

    // Show loading overlay
    showLoading();

    const minLoadTime = new Promise(resolve => setTimeout(resolve, 900));

    try {
        // Send request to backend
        const fetchPromise = fetch('/api/calculate', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({
                name1: name1,
                name2: name2,
            }),
        });

        // Wait for both fetch and minimum playful loading animation
        const [response] = await Promise.all([fetchPromise, minLoadTime]);

        // Parse response
        const data = await response.json();

        if (!response.ok) {
            // Backend validation error
            const errorMsg = data.error || 'Something went wrong. Please try again! 💔';
            showError(errorMsg);
            submitBtn.disabled = false;
            hideLoading();
            return;
        }

        // Hide loading overlay
        hideLoading();

        // Display result
        resultMessage.textContent = data.message;
        percentageNumber.textContent = '0';

        // Transition to result section
        formSection.style.display = 'none';
        resultSection.style.display = 'block';

        // Animate the percentage
        await animatePercentage(data.percentage);

    } catch (error) {
        hideLoading();
        showError('Oops! Cosmic connection failed. Please try again. 📡💔');
        console.error('Fetch error:', error);
        submitBtn.disabled = false;
    }
}

/**
 * Reset the form and return to input screen.
 */
function handleReset() {
    // Clear inputs
    name1Input.value = '';
    name2Input.value = '';
    clearError();

    // Reset button state
    submitBtn.disabled = false;

    // Transition back to form
    resultSection.style.display = 'none';
    formSection.style.display = 'block';

    // Focus on the first input
    name1Input.focus();
}

// Initialize: focus on first input
window.addEventListener('load', () => {
    if (name1Input) {
        name1Input.focus();
    }
});
