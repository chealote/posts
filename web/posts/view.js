const HEADER = document.getElementById("header");
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
  TITLE_ITEM.innerHTML = postContent.title;
  POST_ITEM.value = postContent.contents;
  console.log(POST_ITEM);

  console.log("the post content:", postContent);
}

fetchPostContent();
