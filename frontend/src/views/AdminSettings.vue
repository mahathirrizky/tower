<template>
  <div>
    <h2 class="text-2xl font-bold mb-4">Admin Settings</h2>
    <Card class="w-full md:w-25rem">
      <template #title>Change Password</template>
      <template #content>
        <Form v-slot="$form" :initialValues :resolver @submit="onFormSubmit" class="flex flex-col gap-4">
          <div class="flex flex-col gap-1">
            <FloatLabel variant="on">
              <Password toggleMask id="currentPassword" name="currentPassword" :class="{ 'p-invalid': $form.currentPassword?.invalid }" :feedback="false" fluid />
              <label for="currentPassword">Current Password</label>
            </FloatLabel>
            <Message v-if="$form.currentPassword?.invalid" severity="error" size="small" variant="simple">{{ $form.currentPassword.error.message }}</Message>
          </div>

          <div class="flex flex-col gap-1">
            <FloatLabel variant="on">
              <Password toggleMask id="newPassword" name="newPassword" :feedback="false" :class="{ 'p-invalid': $form.newPassword?.invalid }" fluid />
              <label for="newPassword">New Password</label>
            </FloatLabel>
            <Message v-if="$form.newPassword?.invalid" severity="error" size="small" variant="simple">
                <ul class="my-0 px-4 flex flex-col gap-1">
                    <li v-for="(error, index) of $form.newPassword.errors" :key="index">{{ error.message }}</li>
                </ul>
            </Message>
          </div>

          <div class="flex flex-col gap-1">
            <FloatLabel variant="on">
              <Password toggleMask id="confirmNewPassword" name="confirmNewPassword" :feedback="false" :class="{ 'p-invalid': $form.confirmNewPassword?.invalid }" fluid />
              <label for="confirmNewPassword">Confirm New Password</label>
            </FloatLabel>
            <Message v-if="$form.confirmNewPassword?.invalid" severity="error" size="small" variant="simple">{{ $form.confirmNewPassword.error.message }}</Message>
          </div>

          <Button type="submit" label="Change Password" icon="pi pi-check" class="w-full mt-2" />
        </Form>
      </template>
    </Card>
  </div>
</template>

<script setup>
import { ref } from 'vue';
import { useAuthStore } from '../stores/auth';
import { useToast } from 'primevue/usetoast';
import { z } from 'zod';
import { zodResolver } from '@primevue/forms/resolvers/zod';
import Card from 'primevue/card';
import Button from 'primevue/button';
import FloatLabel from 'primevue/floatlabel';
import { Form } from '@primevue/forms';
import Message from 'primevue/message';
import Password from 'primevue/password';

const authStore = useAuthStore();
const toast = useToast();

const initialValues = ref({
    currentPassword: '',
    newPassword: '',
    confirmNewPassword: ''
});

const resolver = zodResolver(
  z.object({
    currentPassword: z.string().min(1, { message: 'Current password is required.' }),
    newPassword: z
        .string()
        .min(6, { message: 'New password must be at least 6 characters.' })
        .refine((value) => /[a-z]/.test(value), {
            message: 'Must have a lowercase letter.'
        })
        .refine((value) => /[A-Z]/.test(value), {
            message: 'Must have an uppercase letter.'
        })
        .refine((value) => /\d/.test(value), {
            message: 'Must have a number.'
        }),
    confirmNewPassword: z.string().min(1, { message: 'Confirm new password is required.' })
  }).refine((data) => data.newPassword === data.confirmNewPassword, {
    message: 'New passwords do not match.',
    path: ['confirmNewPassword'],
  })
);

const onFormSubmit = async (e) => {
  // e.values: An object containing the current values of all form fields.
  // e.valid: A boolean that indicates whether the form is valid or not.
  if (e.valid) {
    try {
      await authStore.changePassword(e.values.currentPassword, e.values.newPassword);
      toast.add({ severity: 'success', summary: 'Success', detail: 'Password changed successfully', life: 3000 });
    } catch (error) {
      toast.add({ severity: 'error', summary: 'Error', detail: error.message || 'Failed to change password', life: 3000 });
    }
  }
};
</script>