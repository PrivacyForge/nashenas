<template>
  <div class="grid grid-cols-1 gap-4 m-4 bg-[#ffffff] p-5 rounded-lg">
    <p>Write something for Yasha...</p>
    <textarea
      class="textarea textarea-bordered w-full focus:outline-[#119af5]"
      placeholder="Write..."
    ></textarea>
    <button class="btn bg-[#119af5] text-white btn-block">
      Send to @{{ $route.params.username }}
    </button>
  </div>
</template>
