document.addEventListener('DOMContentLoaded', () => {
  const btn = document.querySelector('#btnPosts');
  btn.addEventListener('click', getPosts);

  function getPosts() {
    const div = document.querySelector('#posts');
    const url = '/api/articles';
    fetch(url)
      .then(res => {
        return res.text();
      })
      .then(data => {
        div.innerHTML = data;
      });
  }
});
