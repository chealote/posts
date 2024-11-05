const URL_PARAMS = new URLSearchParams(window.location.search);
const TITLE_PARAM = URL_PARAMS.get("title");
const TITLE_ITEM = document.getElementById("title");
const POST_ITEM = document.getElementById("post");

TITLE_ITEM.value = TITLE_PARAM;

async function fetchPostContent() {
  const response = await fetch(`${BASE_URL}/posts/${TITLE_PARAM}`, {
    headers: {
      "Authorization": sessionStorage.getItem("token"),
    },
  });

  console.log(response);
  if (!response.ok) {
    console.error("failed getting post content");
    return;
  }

  // a json string is a valid json?
  const postContent = await response.json();

  POST_ITEM.value = postContent;
}

fetchPostContent();
