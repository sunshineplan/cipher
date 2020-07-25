const BootstrapButtons = Swal.mixin({
  customClass: {
    confirmButton: 'swal btn btn-primary'
  },
  buttonsStyling: false
});

function Copy() {
  var textarea = $('#' + $("input[name='text']:checked").val());
  if (textarea.val().trim() !== '')
    navigator.clipboard.writeText(textarea.val())
      .then(() => BootstrapButtons.fire('Success', 'Text has been copied to clipboard.', 'success'))
      .catch(() => BootstrapButtons.fire('Error', 'Unable to copy to clipboard.', 'error'));
};

function Clear() {
  $('#unencrypted').val('');
  $('#encrypted').val('');
};

function waiting(on = true) {
  if (on == true) {
    $('.btn').attr('disabled', true);
    $('body').addClass('wait');
    $('.navbar').css('pointer-events', 'none');
    $('.row').css('pointer-events', 'none');
  } else {
    $('.btn').attr('disabled', false);
    $('body').removeClass('wait');
    $('.navbar').css('pointer-events', 'auto');
    $('.row').css('pointer-events', 'auto');
  };
};

function doEncrypt() {
  if ($('#unencrypted').val() == '') {
    BootstrapButtons.fire('Error', 'Empty unencrypted text!', 'error');
    return;
  };
  var key = $('#key').val().trim();
  var promise = new Promise(function (resolve) { if (key == '') resolve(false); else resolve(true) });
  promise.then(haskey => {
    if (!haskey)
      return Swal.fire({
        title: 'Warning!',
        text: 'No key provided, only encode using base64.',
        icon: 'warning',
        confirmButtonText: 'Continue',
        showCancelButton: true,
        customClass: {
          confirmButton: 'swal btn btn-primary',
          cancelButton: 'swal btn btn-danger'
        },
        buttonsStyling: false
      });
    else return { value: true }
  }).then(confirm => {
    if (confirm.value)
      if ($('#online').prop('checked')) {
        waiting();
        $.post('do', {
          mode: 'encrypt',
          key: key,
          content: $('#unencrypted').val()
        }, data => {
          if (data.result != null) {
            $('#encrypted').val(data.result);
            $('textarea').scrollTop(0);
          } else {
            BootstrapButtons.fire('Error', 'Unknow error!', 'error');
          };
        }, 'json')
          .fail(() => BootstrapButtons.fire('Error', 'Network error!', 'error'))
          .always(() => waiting(false));
      } else
        try {
          waiting();
          encrypt();
          $('textarea').scrollTop(0);
        } catch (e) {
          BootstrapButtons.fire('Error', e.message, 'error');
        } finally {
          waiting(false);
        };
  });
};

function doDecrypt() {
  if ($('#encrypted').val() == '') {
    BootstrapButtons.fire('Error', 'Empty encrypted text!', 'error');
    return false;
  };
  if ($('#online').prop('checked')) {
    waiting();
    $.post('do', {
      mode: 'decrypt',
      key: $('#key').val().trim(),
      content: $('#encrypted').val()
    }, data => {
      if (data.result != null) {
        $('#unencrypted').val(data.result);
        $('textarea').scrollTop(0);
      } else BootstrapButtons.fire('Error', 'Incorrect key or malformed encrypted text!', 'error');
    }, 'json')
      .fail(() => BootstrapButtons.fire('Error', 'Network error!', 'error'))
      .always(() => waiting(false));
  } else {
    try {
      waiting();
      decrypt();
      $('textarea').scrollTop(0);
    } catch (e) {
      BootstrapButtons.fire('Error', 'Incorrect key or malformed encrypted text!<br><br>' + e.message, 'error');
    } finally {
      waiting(false);
    };
  };
};

concat = sjcl.bitArray.concat;
base64 = sjcl.codec.base64

function encrypt() {
  var key = $('#key').val().trim();
  var content = $('#unencrypted').val();
  if (key == '') {
    $('#encrypted').val(btoa(unescape(encodeURIComponent(content))).replace(/=/g, ''));
    return;
  };
  sjcl.misc.pa = {};
  var plaintext = zlib(content);
  var data = sjcl.json.ja(key, plaintext.content);
  $('#encrypted').val(base64.fromBits(concat(concat(data.salt, data.iv), concat(data.ct, plaintext.compression))).replace(/=/g, ''));
};

function decrypt() {
  var key = $('#key').val().trim();
  if (key == '') {
    $('#unencrypted').val(decodeURIComponent(escape(atob($('#encrypted').val()))));
    return;
  };
  var cipher = BitsToUint8Array(base64.toBits($('#encrypted').val()));
  var data = {};
  data.salt = Uint8ArrayToBits(cipher.slice(0, 8));
  data.iv = Uint8ArrayToBits(cipher.slice(8, 24));
  data.ct = Uint8ArrayToBits(cipher.slice(24, cipher.length - 1));
  if (new TextDecoder().decode(cipher.slice(cipher.length - 1, cipher.length)) == 1)
    $('#unencrypted').val(zlib(sjcl.json.ia(key, data, { raw: 1 })));
  else $('#unencrypted').val(sjcl.json.ia(key, data));
};

function zlib(obj) {
  if (typeof obj == 'string') {
    var uint8array = new TextEncoder().encode(obj);
    var deflate = pako.deflate(uint8array);
    if (uint8array.length > deflate.length)
      return { content: Uint8ArrayToBits(deflate), compression: sjcl.codec.utf8String.toBits(1) };
    else return { content: Uint8ArrayToBits(uint8array), compression: sjcl.codec.utf8String.toBits(0) };
  } else return pako.inflate(BitsToUint8Array(obj), { to: 'string' })
};

function Uint8ArrayToBits(uint8array) {
  var hex = '';
  for (var i = 0; i < uint8array.length; i++)
    hex += (uint8array[i] + 0xF00).toString(16).substr(1);
  return sjcl.codec.hex.toBits(hex);
};

function BitsToUint8Array(bits) {
  var array = [], hex = sjcl.codec.hex.fromBits(bits);
  for (var i = 0; i < hex.length; i += 2)
    array.push(parseInt(hex.substr(i, 2), 16));
  return new Uint8Array(array);
};
