
// Show notifications
function  notification(icon, title){
    if(icon === "success"){
        Swal.fire({
        title: title,
        icon: icon,
        draggable: true
        });
    }else if(icon === "error"){
        Swal.fire({
        icon: icon,
        title:  title,
        });
    }
}

window.notification = notification;
