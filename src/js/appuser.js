//Fetch all roles here and populate html element 
fetch('/api/roles')
  .then(res => res.json())
  .then(data => {
      AppState.roles = data;

      const select = document.getElementById('roleid');
      if (!select) return;

      // clear existing options
      select.innerHTML = '<option value="">Select role</option>';

      data.forEach(r => {
          const option = document.createElement('option');
          option.value = r.roleid;
          option.textContent = r.rolename; 
          select.appendChild(option);
      });
  })
  .catch(err => console.error('Failed to load roles:', err));


