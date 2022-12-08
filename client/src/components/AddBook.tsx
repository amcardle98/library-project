import React, { useState } from 'react'
import { useForm } from '@mantine/form'
import { Modal, Group, Button, Card, CardSection } from '@mantine/core'
import { Input, Button as MatButton } from '@mui/material'
import { ENDPOINT } from '../App'
import { KeyedMutator } from 'swr'

function AddBook({ mutate }: { mutate: KeyedMutator<any> }){
    const [open, setOpen] = useState(false)

    const form = useForm({
        initialValues: { title: ''},
    })

async function createBook(values: {title: string}){
    const updated = await fetch(`${ENDPOINT}/api/create`, {
        method: 'POST',
        headers: {  'Content-Type': 'application/json' },
        body: JSON.stringify(values)
    }).then((res) => res.json().catch(() => null));

    mutate(updated)
    form.reset()
    setOpen(false)
}

    return (
        <>
            <Modal opened={open} onClose={() => setOpen(false)} title='Add Book'>
                     <form onSubmit={form.onSubmit(createBook)}>
                        <Input placeholder='Enter book title' required {...form.getInputProps("title")}></Input>
                        <MatButton type="submit">Create Book</MatButton>
                     </form>
            </Modal>
            <Group position="center">
                <Button onClick={() => setOpen(true)}>Add Book</Button>
            </Group>
        </>
    );
}

export default AddBook