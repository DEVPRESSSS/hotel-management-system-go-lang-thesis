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


//Fetch all floor here and populate html element 
fetch('/api/floors')
  .then(res => res.json())
  .then(floors => {
      AppState.floors = floors;

      const select = document.getElementById('floorid');
      if (!select) return;

      // clear existing options
      select.innerHTML = '<option value="">Select floor</option>';

      floors.forEach(r => {
          const option = document.createElement('option');
          option.value = r.floorid;
          option.textContent = r.floorname; 
          select.appendChild(option);
      });
  })
.catch(err => console.error('Failed to load roles:', err));

//Fetch all floor here and populate html element 
fetch('/api/roomtypes')
  .then(res => res.json())
  .then(rt => {
      AppState.rt = rt;

      const select = document.getElementById('roomtypeid');
      if (!select) return;

      // clear existing options
      select.innerHTML = '<option value="">Select room type</option>';

      rt.forEach(r => {
          const option = document.createElement('option');
          option.value = r.roomtypeid;
          option.textContent = r.roomtypename; 
          select.appendChild(option);
      });
  })
.catch(err => console.error('Failed to load roles:', err));

//Fetch access
fetch('api/access')
  .then(res => res.json())
  .then(access => {
      AppState.access = access;

      const select = document.getElementById('accessid');
      if (!select) return;

      select.innerHTML = '<option value="">Select floor</option>';

      access.forEach(r => {
          const option = document.createElement('option');
          option.value = r.accessid
          option.textContent = r.accessname; 
          select.appendChild(option);
      });
  })
.catch(err => console.error('Failed to load roles:', err));

//Fetch amenity
fetch('api/aminities')
  .then(res => res.json())
  .then(aminities => {
      AppState.aminities = aminities;

      const select = document.getElementById('aminityid');
      if (!select) return;

      select.innerHTML = '<option value="">Select Aminity</option>';

      aminities.forEach(r => {
          const option = document.createElement('option');
          option.value = r.amenityid
          option.textContent = r.amenityname; 
          select.appendChild(option);
      });
  })
.catch(err => console.error('Failed to load aminity:', err));


//Fetch rooms
fetch('api/rooms')
  .then(res => res.json())
  .then(rooms => {
      AppState.rooms = rooms;

      const select = document.getElementById('roomid');
      if (!select) return;

      select.innerHTML = '<option value="">Select Room</option>';

      rooms.forEach(r => {
          const option = document.createElement('option');
          option.value = r.roomid
          option.textContent = r.roomnumber; 
          select.appendChild(option);
      });
  })
.catch(err => console.error('Failed to load aminity:', err));