const BASE_URL = 'http://127.0.0.1:8080';
const SIGNUP_FORM = document.getElementById('signup-form');

async function signup(event) {
  event.preventDefault();

  const username = document.getElementById('username').value;
  const password = document.getElementById('password').value;

  try {
    const response = await fetch(`${BASE_URL}/signup`, {
      method: 'POST',
      credentials: 'include',
      headers: {
        'Content-Type': 'application/json',
        'Access-Control-Allow-Origin': 'index.html',
      },
      body: JSON.stringify({
        'name': username,
        'password': password,
      })
    });

    if (response.ok) {
      window.location.replace('../index.html');
    } else {
      const errorText = document.getElementById('error-text');
      errorText.textContent = 'Signup failed. Please check your credentials.';
    }
  } catch (err) {
    const errorText = document.getElementById('error-text');
    errorText.textContent = 'Signup failed. Please try again later.';
  }
}
