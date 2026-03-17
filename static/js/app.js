function Prompt() {
  let toast = function (c) {
    const { msg = "", icon = "success", position = "top-end" } = c;

    const Toast = Swal.mixin({
      toast: true,
      title: msg,
      position: position,
      icon: icon,
      showConfirmButton: false,
      timer: 3000,
      timerProgressBar: true,
      didOpen: (toast) => {
        toast.addEventListener("mouseenter", Swal.stopTimer);
        toast.addEventListener("mouseleave", Swal.resumeTimer);
      },
    });

    Toast.fire({});
  };

  let success = function (c) {
    const { msg = "", title = "", footer = "" } = c;

    Swal.fire({
      icon: "success",
      title: title,
      text: msg,
      footer: footer,
    });
  };

  let error = function (c) {
    const { msg = "", title = "", footer = "" } = c;

    Swal.fire({
      icon: "error",
      title: title,
      text: msg,
      footer: footer,
    });
  };

  async function custom(c) {
    const { icon = "", msg = "", title = "", showConfirmButton = true } = c;

    const { value: result } = await Swal.fire({
      icon: icon,
      title: title,
      html: msg,
      backdrop: false,
      focusConfirm: false,
      showCancelButton: true,
      showConfirmButton: showConfirmButton,

      willOpen: () => {
        if (c.willOpen !== undefined) {
          c.willOpen();
        }
      },

      didOpen: () => {
        if (c.didOpen !== undefined) {
          c.didOpen();
        }
      },

      preConfirm: () => {
        return [
          document.getElementById("start").value,
          document.getElementById("end").value,
        ];
      },
    });

    if (result) {
      if (result.dismiss !== Swal.DismissReason.cancel) {
        if (result.value !== "") {
          if (c.callback !== undefined) {
            c.callback(result);
          } else {
            c.callback(false);
          }
        } else {
          c.callback(false);
        }
      }
    }
  }

  return {
    toast: toast,
    success: success,
    error: error,
    custom: custom,
  };
}

async function checkAvailabilityForm(result, token, id) {
  let form = document.getElementById("check-availability-form");
  let formData = new FormData(form);
  // append the CSRF token from base layout (meta) to the form data
  formData.append("csrf_token", `${token}`);
  formData.append("room_id", `${id}`);

  const res = await fetch("/search-availability-json", {
    method: "POST",
    body: formData,
  });
  const data = await res.json();
  if (data.ok) {
    attention.custom({
      icon: "success",
      showConfirmButton: false,
      msg:
        "<p>Room is available</p>" +
        '<p><a href="/book-room?id=' +
        data.room_id +
        "&s=" +
        data.start_date +
        "&e=" +
        data.end_date +
        '" class="btn btn-primary">Book now!</a></p>',
    });
  } else {
    attention.error({
      msg: "No avalability",
    });
  }
}
