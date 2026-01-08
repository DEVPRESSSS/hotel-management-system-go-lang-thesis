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
    headerTitle.innerText = "Create room aminity";
    btnSubmit.innerText = "Create";
    openModal();
}

  //Upsert role
document.getElementById('upsertform').addEventListener('submit', function(e){

    e.preventDefault();
    //Generate unique uid
    let uid = "";
    //Get the input in role textbox
    const roomid = document.getElementById('roomid').value;
    // const roomid = document.getElementById('roomid').value;
    if( id === ""){
        uid = uuidv4();
    }else{
        uid = id;
    }

    const formData ={
        aminityid: uid,
        aminityname: aminityname,

    };
    if(id === "" || id === null){
        fetch('/api/createaminity', {
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

                //Show notication success
                notification("success", data.success);

            })
            .catch(err=>{

                notification("error", err.message);

            });
    }else{
           //Update user
            const updateFormData = {
                  aminityid: uid,
                  aminityname: aminityname,       
            };
            fetch(`/api/updateroomaminity/${id}`, {
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

//Fetch all room aminities
fetch('/api/roomaminities')
    .then(response => response.json())
        .then(roomAminities => {

            const tbody = document.getElementById('users-body');

            roomAminities.forEach(aminity => {
                tbody.innerHTML += `
                
                    <tr>
                        <td>${aminity.Room.roomnumber}</td>
                        <td>${aminity.Amenity.aminityname}</td>                  
                        <td>
                            <button class="update-btn text-blue-600 hover:text-blue-800 font-medium"
                                data-roomid = "${aminity.RoomId}"
                                >Edit
                            
                            </button>
                            <button 
                                class="delete-btn text-red-600 hover:text-red-800 font-medium"
                                data-roomid="${aminity.RoomId}">
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
    headerTitle.innerText = "Update aminity";  
        
    
    fetch(`/api/roomaminity/${id}`)
        .then(async res => {
            const data = await res.json();
            if (!res.ok) throw new Error(data.error || "Request failed");
            return data;
        })
        .then(data => {
            //notification("success", data.success || "Operation successful");           
            const aminityname= data.success;
            document.getElementById('aminity').value = aminityname.aminityname;

        })
        .catch(err => {
            console.log(err);
            alert(`${err}`);
        });

    openModal();

}); 


//Delete user function
document.getElementById('users-body').addEventListener('click', function (e) {
    if (!e.target.classList.contains('delete-btn')) return;

    const aminityid = e.target.dataset.roomid;
    if (!aminityid) return;

    Swal.fire({
        title: "Are you sure?",
        text: "You won't be able to revert this!",
        icon: "warning",
        showCancelButton: true,
        confirmButtonColor: "#3085d6",
        cancelButtonColor: "#d33",
        confirmButtonText: "Yes, delete it!"
    }).then((result) => {
        if (result.isConfirmed) {
            fetch(`/api/deleteroomaminity/${aminityid}`, {
            method: 'DELETE'

            
        },  
            Swal.fire({
            title: "Deleted!",
            text: "Your file has been deleted.",
            icon: "success"
            })
        )
        .then(res => {
            if (!res.ok) throw new Error('Delete failed');
            e.target.closest('tr').remove(); 
            
        })
        .catch(err => {
            console.log(err);
            alert('Error deleting user');
        });

      
    }
    });

   
});