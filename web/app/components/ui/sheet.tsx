import type { ReactNode } from "react";
import { ModalOverlay, Modal, Dialog, Heading } from "react-aria-components";

type Props = {
  isOpen: boolean;
  onOpenChange: (open: boolean) => void;
  title?: string;
  children: ReactNode;
};

const overlayCls =
  "fixed inset-0 z-40 flex items-end justify-center bg-slate-900/40 " +
  "data-[entering]:animate-[piamap-fade-in_180ms_ease-out] " +
  "data-[exiting]:animate-[piamap-fade-out_140ms_ease-in]";

const modalCls =
  "w-full max-w-[390px] rounded-t-2xl bg-white shadow-xl outline-none " +
  "data-[entering]:animate-[piamap-sheet-in_220ms_cubic-bezier(0.32,0.72,0,1)] " +
  "data-[exiting]:animate-[piamap-sheet-out_180ms_cubic-bezier(0.32,0.72,0,1)]";

export function Sheet({ isOpen, onOpenChange, title, children }: Props) {
  return (
    <ModalOverlay isOpen={isOpen} onOpenChange={onOpenChange} isDismissable className={overlayCls}>
      <Modal className={modalCls}>
        <Dialog className="relative outline-none">
          <div className="flex justify-center pb-1 pt-2">
            <span className="block h-1.5 w-10 rounded-full bg-slate-300" aria-hidden />
          </div>
          {title ? (
            <Heading slot="title" className="px-5 pb-2 text-base font-bold text-slate-900">
              {title}
            </Heading>
          ) : null}
          <div className="max-h-[70vh] overflow-y-auto px-5 pb-6">{children}</div>
        </Dialog>
      </Modal>
    </ModalOverlay>
  );
}

// 中央寄せの通常モーダル (LP の Sign up prompt 等で使用)。
export function CenterDialog({ isOpen, onOpenChange, title, children }: Props) {
  return (
    <ModalOverlay
      isOpen={isOpen}
      onOpenChange={onOpenChange}
      isDismissable
      className={
        "fixed inset-0 z-40 flex items-center justify-center bg-slate-900/40 px-6 " +
        "data-[entering]:animate-[piamap-fade-in_180ms_ease-out] " +
        "data-[exiting]:animate-[piamap-fade-out_140ms_ease-in]"
      }
    >
      <Modal
        className={
          "w-full max-w-[340px] rounded-2xl bg-white shadow-xl outline-none " +
          "data-[entering]:animate-[piamap-zoom-in_180ms_cubic-bezier(0.32,0.72,0,1)] " +
          "data-[exiting]:animate-[piamap-fade-out_140ms_ease-in]"
        }
      >
        <Dialog className="p-5 outline-none">
          {title ? (
            <Heading slot="title" className="text-base font-bold text-slate-900">
              {title}
            </Heading>
          ) : null}
          <div className="mt-2">{children}</div>
        </Dialog>
      </Modal>
    </ModalOverlay>
  );
}
