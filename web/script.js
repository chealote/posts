const CONTENT_DIV = document.getElementById('content');

async function fetchContent(path) {
  const response = await fetch(`${BASE_URL}${path}`, {
    method: 'GET',
    headers: {
      "Authorization": sessionStorage.getItem("token"),
    },
  });

  return new Promise((resolve, reject) => {
    if (response.ok) {
      return resolve(response.text());
    }
    reject('auth failed?');
  });
}

function loadContent(path) {
  // handle error or something
  if (!path || path === '') {
    path = "/";
  }
  CONTENT_DIV.innerHTML = "loading...";
  const content = fetchContent(path)
    .then(content => {
      CONTENT_DIV.innerHTML = content;
    },
      err => {
        console.error("some error:", err);
      });
}

async function listPosts() {
  const response = await fetch(`${BASE_URL}/posts`, {
    headers: {
      "Authorization": sessionStorage.getItem("token"),
    }
  });

  if (! response.ok) {
    return;
  }
  console.log("response from listPosts:", response);

  const postListDiv = document.getElementById("post-list");
  const postList = await response.json();
  for (const post of postList) {
    const link = document.createElement("a");
    link.href = "#";
    link.innerHTML = post.title;

    postListDiv.appendChild(link);
  }
}

listPosts();
