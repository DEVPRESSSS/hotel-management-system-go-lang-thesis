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
    headerTitle.innerText = "Create facility";
    btnSubmit.innerText = "Create";
    openModal();
}



  //Upsert role
document.getElementById('upsertform').addEventListener('submit', function(e){

    e.preventDefault();
    //Generate unique uid
    let uid = "";
    //Get the input in role textbox
    
    const facilityName = document.getElementById('facilityname').value;
    const dateInput = document.getElementById("maintenance_date").value;
    const isoDate = new Date(dateInput + "T00:00:00").toISOString();

      


    if( id === ""){
        uid = uuidv4();
    }else{
        uid = id;
    }

    const formData ={
        facility_id: uid,
        facility_name: facilityName,
        maintenance_date: isoDate,

    };
    if(id === "" || id === null){
        fetch('/api/createfacility', {
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
                facility_id: uid,
                facility_name: facilityName,
                maintenance_date: isoDate,       
            };
            fetch(`/api/updatefacility/${id}`, {
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

//Fetch all roles
fetch('/api/facility')
    .then(response => response.json())
        .then(facilities => {

            const tbody = document.getElementById('users-body');

            facilities.forEach(facility => {
                tbody.innerHTML += `
                
                    <tr>
                        <td>${facility.facility_id}</td>
                        <td>${facility.facility_name}</td>                  
                        <td>${facility.status}</td>                  
                        <td>${new Date(facility.maintenance_date).toLocaleDateString()}</td>                  
                        <td>${new Date(facility.created_at).toLocaleDateString()}</td>                  
                        <td>
                            <button class="update-btn text-blue-600 hover:text-blue-800 font-medium"
                                data-facility_id = "${facility.facility_id}"
                                >Edit
                            
                            </button>
                            <button 
                                class="delete-btn text-red-600 hover:text-red-800 font-medium"
                                data-facility_id="${facility.facility_id}">
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

    id = e.target.dataset.facility_id;
    if (!id) return;

    const facility_id = document.getElementById('input-id');
    facility_id.value = id;

    btnSubmit.innerText = "Update";
    headerTitle.innerText = "Update user";  
        
    
    fetch(`/api/facility/${id}`)
        .then(async res => {
            const data = await res.json();
            if (!res.ok) throw new Error(data.error || "Request failed");
            return data;
        })
        .then(facility => {
            //notification("success", data.success || "Operation successful");           
            const fc= facility.success;
            document.getElementById('facilityname').value = fc.facility_name;
            const date = new Date(fc.maintenance_date);
            const formattedDate = date.toISOString().split('T')[0];

            document.getElementById('maintenance_date').value = formattedDate;
        })
        .catch(err => {
            console.log(err);
                alert(`${err}`);
            //notification("error", err.err);
        });

    openModal();

}); 


//Delete user function
document.getElementById('users-body').addEventListener('click', function (e) {
    if (!e.target.classList.contains('delete-btn')) return;

    const facility_id = e.target.dataset.facility_id;
    if (!facility_id) return;

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
            fetch(`/api/deletefacility/${facility_id}`, {
            method: 'DELETE'
        })
        .then(res => {
            if (!res.ok) throw new Error('Delete failed');
            e.target.closest('tr').remove(); 
            
        })
        .catch(err => {
            console.log(err);
            alert('Error deleting user');
        });

        Swal.fire({
        title: "Deleted!",
        text: "Your file has been deleted.",
        icon: "success"
        });
    }
    });

   
});