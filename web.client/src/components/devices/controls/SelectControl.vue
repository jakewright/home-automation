<template>
    <div>
        <select :value="value" @change="$emit('input', $event.target.value)">
            <option disabled value="">Please select one</option>
            <option v-for="option in options" :value="option.value" :key="option.value">
                {{ option.text }}
            </option>
        </select>
    </div>
</template>

<script>
    export default {
        name: 'SelectControl',

        props: {
            value: {
                type: String,
                required: true,
            },

            options: {
                type: Array,
                required: true,
                validator: function (options) {
                    for (let option of options) {
                        if (!option.hasOwnProperty('value')) return false;
                        if (!option.hasOwnProperty('text')) return false;
                    }

                    return true;
                }
            }
        }
    }
</script>
