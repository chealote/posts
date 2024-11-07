const TITLE_INPUT = document.getElementById("title");
const POST_INPUT = document.getElementById("post");

let INVALID_FORM = false;

function ensureNotEmpty(input) {
  console.log(input.value);
  if (!input.value || input.value === "") {
    input.classList.add("invalid");
    input.invalid = true;
  } else {
    input.classList.remove("invalid");
    input.invalid = false;
  }
}

async function create() {
  ensureNotEmpty(TITLE_INPUT);
  ensureNotEmpty(POST_INPUT);

  const isFormValid = !TITLE_INPUT.invalid && !POST_INPUT.invalid;
  if (!isFormValid) {
    return;
  }

  const response = await fetch(`${BASE_URL}/posts`, {
    method: "POST",
    headers: {
      "Authorization": sessionStorage.getItem("token"),
    },
    body: JSON.stringify({
      title: TITLE_INPUT.value,
      post: POST_INPUT.value,
    }),
  });
  if (response.ok) {
    authRedirect("/");
  }
}
