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
    headerTitle.innerText = "Create service";
    btnSubmit.innerText = "Create";
    openModal();
}



  //Upsert role
document.getElementById('upsertform').addEventListener('submit', function(e){

    e.preventDefault();
    //Generate unique uid
    let uid = "";
    //Get the input in role textbox
    const serviceName = document.getElementById('servicename').value;
    const startTimeInput = document.getElementById('start-time').value;
    const endTimeInput = document.getElementById('end-time').value;
   
    if( id === ""){
        uid = uuidv4();
    }else{
        uid = id;
    }

    const formData ={
        serviceid: uid,
        servicename: serviceName,
        start_time: startTimeInput,
        end_time: endTimeInput,

    };
    if(id === "" || id === null){
        fetch('/api/createservice', {
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
                service_id: uid,
                service_name: serviceName,
                start_time: startTimeInput,
                end_time: endTimeInput,     
            };
            fetch(`/api/updateservice/${id}`, {
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
fetch('/api/services')
    .then(response => response.json())
        .then(services => {

            const tbody = document.getElementById('users-body');

            services.forEach(service => {
                tbody.innerHTML += `
                
                    <tr>
                        <td>${service.serviceid}</td>
                        <td>${service.servicename}</td>                  
                        <td>${service.start_time}</td> 
                        <td>${service.end_time}</td>                  
                        <td>${service.created_at}</td>                  
                        <td>
                            <button class="update-btn text-blue-600 hover:text-blue-800 font-medium"
                                data-serviceid = "${service.serviceid}"
                                >Edit
                            
                            </button>
                            <button 
                                class="delete-btn text-red-600 hover:text-red-800 font-medium"
                                data-serviceid="${service.serviceid}">
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

    id = e.target.dataset.serviceid;
    if (!id) return;

    const serviceid = document.getElementById('input-id');
    serviceid.value = id;

    btnSubmit.innerText = "Update";
    headerTitle.innerText = "Update user";  
        
    
    fetch(`/api/service/${id}`)
        .then(async res => {
            const data = await res.json();
            if (!res.ok) throw new Error(data.error || "Request failed");
            return data;
        })
        .then(service => {
            //notification("success", data.success || "Operation successful");           
            const sv= service.success;
            document.getElementById('servicename').value = sv.servicename;
            const st = new Date(sv.start_time);
            const formattedDateSt = st.toISOString().split('T')[0];
            document.getElementById('start-time').value = formattedDateSt;

            const et = new Date(sv.end_time);
            const formattedDateEt = st.toISOString().split('T')[0];
            document.getElementById('end-time').value = formattedDateEt;
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

    const serviceid = e.target.dataset.serviceid;
    if (!serviceid) return;

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
            fetch(`/api/deleteservice/${serviceid}`, {
            method: 'DELETE'
        })
        .then(res => {
            if (!res.ok) throw new Error('Delete failed');
            e.target.closest('tr').remove(); 
            
        })
        .catch(err => {
            console.log(err);
            alert('Error deleting user');
            return;
        });

       
    }
     Swal.fire({
        title: "Deleted!",
        text: "Your file has been deleted.",
        icon: "success"
        });
    });

   
});