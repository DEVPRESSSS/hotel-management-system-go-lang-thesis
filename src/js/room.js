//Access the header title to modify it later
const headerTitle = document.getElementById('header-action');
//Access the header text to modify its text inside
const btnSubmit = document.getElementById('btn-submit');
//Stors user id globally
let id = "";
function openModal() {
    document.getElementById('userModal').classList.remove('hidden');
    document.getElementById('userModal').classList.add('flex');
}

function closeModal() {
        document.getElementById('userModal').classList.add('hidden');
        document.getElementById('userModal').classList.remove('flex');
}
function createModal(){
    headerTitle.innerText = "Create room";
    btnSubmit.innerText = "Create";
    openModal();
}



  //Upsert role
document.getElementById('upsertform').addEventListener('submit', function(e){

    e.preventDefault();

    let uid = "";
    const roomNumber = document.getElementById('roomnumber').value;
    const roomType = document.getElementById('roomtypeid').value;
    const floor = document.getElementById('floorid').value;
    const capacity = document.getElementById('capacity').value;
    const price = document.getElementById('price').value;
   
    if( id === ""){
        uid = uuidv4();
    }else{
        uid = id;
    }

    const formData ={
        roomid: uid,
        roomnumber: roomNumber,
        roomtypeid: roomType,
        floorid: floor,
        capacity: capacity,
        price: price

    };
    if(id === "" || id === null){
        fetch('/api/createroom', {
            method:'POST',
            headers:{
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(formData),
            })
            .then(async response => {
                        const data = await response.json();

                        if (!response.ok) {
                            throw new Error(data.error);
                        }
                        return data;
            
            })
            .then(data=>{
                notification("success", data.success);
            })
            .catch(err=>{

                notification("error", err.message);

            });
    }else{
           //Update user
            const updateFormData = {
                roomid: uid,
                roomnumber: roomNumber,
                roomtypeid: roomType,
                floorid: floor,
                capacity: capacity,
                price: price    
            };
            fetch(`/api/updateroom/${id}`, {
                method: 'PUT',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify(updateFormData)
                })
                .then(async response => {
                    const data = await response.json();

                    if (!response.ok) {
                        throw new Error(data.error);
                    }
                    return data;
                })
                .then(data => {
                    console.log(updateFormData);
                    notification("success", data.success);

                })
                .catch(err => {
                    notification("error", err.error);
            });
    }

    //Close the modal
    closeModal();

});

function uuidv4() {
    return 'xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx'
    .replace(/[xy]/g, function (c) {
        const r = Math.random() * 16 | 0, 
            v = c == 'x' ? r : (r & 0x3 | 0x8);
        return v.toString(16);
    });
}

//Fetch all roles
fetch('/api/rooms')
    .then(response => response.json())
        .then(rooms => {

            const tbody = document.getElementById('users-body');

            rooms.forEach(room => {
                tbody.innerHTML += `
                
                    <tr>
                        <td>${room.roomid}</td>
                        <td>${room.roomnumber}</td>                  
                        <td>${room.RoomType.roomtypename}</td> 
                        <td>${room.Floor.floorname}</td>                  
                        <td>${room.capacity}</td>                  
                        <td>${room.price}</td>                  
                        <td>${room.status}</td>                  
                        <td>${new Date(room.created_at).toLocaleDateString()}</td>                  
                        <td>
                            <button class="update-btn text-blue-600 hover:text-blue-800 font-medium"
                                data-roomid = "${room.roomid}"
                                >Edit
                            
                            </button>
                            <button 
                                class="delete-btn text-red-600 hover:text-red-800 font-medium"
                                data-roomid="${room.roomid}">
                                Delete
                            </button>
                        </td>
                    </tr>
                
                `;               
            });
});

//Edit click 
document.getElementById('users-body').addEventListener('click' , function(e){

    if (!e.target.classList.contains('update-btn')) return;

    id = e.target.dataset.roomid;
    if (!id) return;

    const roomid = document.getElementById('input-id');
    roomid.value = id;

    btnSubmit.innerText = "Update";
    headerTitle.innerText = "Update room";  
        
    
    fetch(`/api/room/${id}`)
        .then(async res => {
            const data = await res.json();
            if (!res.ok) throw new Error(data.error || "Request failed");
            return data;
        })
        .then(room => {
            //notification("success", data.success || "Operation successful");           
            const r= room.success;
            document.getElementById('roomnumber').value= r.roomnumber;
            document.getElementById('roomtypeid').value = r.roomtypeid;
            document.getElementById('floorid').value = r.floorid;
            document.getElementById('capacity').value= r.capacity;
            document.getElementById('price').value= r.price;
   
     
        })
        .catch(err => {
            console.log(err);
            alert(`${err}`);
            //notification("error", err.err);
        });

    openModal();

}); 


document.getElementById('users-body').addEventListener('click', function (e) {
    if (!e.target.classList.contains('delete-btn')) return;

    const roomid = e.target.dataset.roomid;
    if (!roomid) return;

    Swal.fire({
        title: "Are you sure?",
        text: "You won't be able to revert this!",
        icon: "warning",
        showCancelButton: true,
        confirmButtonColor: "#3085d6",
        cancelButtonColor: "#d33",
        confirmButtonText: "Yes, delete it!"
    }).then((result) => {

        if (!result.isConfirmed) return;

        fetch(`/api/deleteroom/${roomid}`, {
            method: 'DELETE'
        })
        .then(res => {
            if (!res.ok) throw new Error('Delete failed');

            Swal.fire({
                title: "Deleted!",
                text: "Your file has been deleted.",
                icon: "success"
            });

            e.target.closest('tr').remove();
        })
        .catch(err => {
            console.error(err);
            Swal.fire("Error", "Failed to delete room", "error");
        });
    });
});
