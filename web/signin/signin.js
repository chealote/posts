const BASE_URL = 'http://127.0.0.1:8080';
const SIGNIN_FORM = document.getElementById('signin-form');

async function signin(event) {
  event.preventDefault();

  const username = document.getElementById('username').value;
  const password = document.getElementById('password').value;

  try {
    const response = await fetch(`${BASE_URL}/signin`, {
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
      // cookie should be set here, check in index if the cookie is set
      window.location.replace('../index.html');
    } else {
      const errorText = document.getElementById('error-text');
      errorText.textContent = 'Signin failed. Please check your credentials.';
    }
  } catch (err) {
    const errorText = document.getElementById('error-text');
    errorText.textContent = 'Signin failed. Please try again later.';
  }
}
