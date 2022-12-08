import React, { useState } from 'react'
import { useForm } from '@mantine/form'
import { Modal, Group, Button, Card, CardSection } from '@mantine/core'
import { Input, Button as MatButton } from '@mui/material'
import { ENDPOINT } from '../App'
import { KeyedMutator } from 'swr'

function FindBook({ mutate }: { mutate: KeyedMutator<any> }){
    const [open, setOpen] = useState(false)

    const form = useForm({
        initialValues: { title: ''},
    })

async function findBookByTitle(values: {title: string}){
    const updated = await fetch(`${ENDPOINT}/api/getByTitle/${values.title}`, {
        method: 'GET',
        headers: {  'Content-Type': 'application/json' }
    }).then((res) => res.json().catch(() => null));

    //find book by title then return book object 

    mutate(updated)
    form.reset()
    setOpen(false)
}

    return (
        <>
            <Modal opened={open} onClose={() => setOpen(false)} title='Find Book'>
                     <form onSubmit={form.onSubmit(findBookByTitle)}>
                        <Input placeholder='Enter book title' required {...form.getInputProps("title")}></Input>
                        <MatButton type="submit">Find Book</MatButton>
                     </form>
            </Modal>
            <Group position="center">
                <Button onClick={() => setOpen(true)}>Find Book</Button>
            </Group>
        </>
    );
}

export default FindBook