$(() => {
  $('select').material_select();
  $('.carousel').carousel();
  $('form').submit(submitForm);
  $('.back').click(hidePair);
  $(window).resize(resizeCarousel);
})

function submitForm(event) {
  event.preventDefault();
  let name = $('select.name-select').val();
  if (!name) {
    $('input').addClass('invalid');
    $('.error').html('Name is Required');
    return
  }
  let animal = $('select.animal-select').val();
  if (!animal) {
    animal = findAnimal(); 
  }
  $('button').hide();
  $('.preloader-wrapper').show();
  $('.error').text('');
  $('.carousel').css('background', 'none')
  $.get(`/pair?name=${name}&animal=${animal}`)
    .done(pair => {
      if (pair.error) {
        $('.error').text(pair.error);
        $('button').show();
        $('.preloader-wrapper').hide();
        $('.carousel').css('background', '#ED5454')
      } else {
        $('button').show();
        $('.preloader-wrapper').hide();
        displayPair(pair);
      }
    })
    .fail(err => {
      console.log(err);
    })
}

function findAnimal() {
  let animal;
  let highestOpacity = 0;
  $('.carousel-item').each(function(index) {
    let opacity = $(this).css('opacity'); 
    if (opacity == 1) {
      animal = $(this).data('animal');
      return
    } else if (opacity > highestOpacity) {
      highestOpacity = opacity;
      animal = $(this).data('animal');
    }
  })
  return animal;
}

function displayPair(pair) {
  $('.form').hide();
  $('.name').text(pair.name);
  $('.card-image img').attr('src',`assets/safari.jpg`);
  $('.pair').show();
  $('.back').show();
}

function hidePair(pair) {
  $('.form').show();
  $('.pair').hide();
  $('.back').hide();
}

function resizeCarousel() {
  $('.carousel').removeClass('initialized');
  $('.carousel').carousel();
}
