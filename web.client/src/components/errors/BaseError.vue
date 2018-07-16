<template>
    <div class="notification is-danger">
        <button class="delete" @click="$emit('close')"></button>
        <slot></slot>
    </div>
</template>
